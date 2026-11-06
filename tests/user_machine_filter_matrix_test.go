package tests

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// filterTree is a company's organization with one machine line at each depth of the tree, so a
// filter at any level has machines just inside and just outside it.
type filterTree struct {
	country, state, city, branch, location, shop, line int64
	lines                                              []int64 // lines[i] leaves the base path at depth i (see newFilterTree)
}

// newFilterTree builds the organization under company's fixture line:
//
//	lines[0] base line                 lines[1] same shop        lines[2] same location
//	lines[3] same branch               lines[4] same city        lines[5] same state
//	lines[6] same country              lines[7] another country
func newFilterTree(t *testing.T, company testCompany, suffix string) filterTree {
	t.Helper()
	tree := filterTree{line: company.lineID}
	err := testDBPool.QueryRow(context.Background(), `SELECT user_country_id, user_state_id, user_city_id, branch_id, user_location_id, user_shop_id
		FROM organization_paths WHERE user_line_id = $1`, company.lineID).
		Scan(&tree.country, &tree.state, &tree.city, &tree.branch, &tree.location, &tree.shop)
	if err != nil {
		t.Fatalf("find the fixture's path: %v", err)
	}
	owner := company.customerID
	addCountry := func(tag string) int64 {
		return insertTestRow(t, "user_country", `INSERT INTO user_country (user_id, name) VALUES ($1, $2) RETURNING id`, owner, "Country"+tag)
	}
	addState := func(country int64, tag string) int64 {
		return insertTestRow(t, "user_state", `INSERT INTO user_state (user_id, user_country_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, country, "State"+tag)
	}
	addCity := func(state int64, tag string) int64 {
		return insertTestRow(t, "user_city", `INSERT INTO user_city (user_id, user_state_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, state, "City"+tag)
	}
	addBranch := func(city int64, tag string) int64 {
		return insertTestRow(t, "branches", `INSERT INTO branches (code, name, user_city_id) VALUES ($1, $2, $3) RETURNING id`, "X"+tag, "Branch"+tag, city)
	}
	addLocation := func(branch int64, tag string) int64 {
		return insertTestRow(t, "user_location", `INSERT INTO user_location (user_id, branch_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, branch, "Location"+tag)
	}
	addShop := func(location int64, tag string) int64 {
		return insertTestRow(t, "user_shop", `INSERT INTO user_shop (user_id, user_location_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, location, "Shop"+tag)
	}
	addLine := func(shop int64, tag string) int64 {
		return insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, shop, "Line"+tag)
	}
	fromShop := func(shop int64, tag string) int64 { return addLine(shop, tag) }
	fromLocation := func(location int64, tag string) int64 { return fromShop(addShop(location, tag), tag) }
	fromBranch := func(branch int64, tag string) int64 { return fromLocation(addLocation(branch, tag), tag) }
	fromCity := func(city int64, tag string) int64 { return fromBranch(addBranch(city, tag), tag) }
	fromState := func(state int64, tag string) int64 { return fromCity(addCity(state, tag), tag) }
	fromCountry := func(country int64, tag string) int64 { return fromState(addState(country, tag), tag) }

	tag := func(depth int) string { return fmt.Sprintf("%d%s", depth, suffix) }
	tree.lines = []int64{
		company.lineID,
		fromShop(tree.shop, tag(1)),
		fromLocation(tree.location, tag(2)),
		fromBranch(tree.branch, tag(3)),
		fromCity(tree.city, tag(4)),
		fromState(tree.state, tag(5)),
		fromCountry(tree.country, tag(6)),
		fromCountry(addCountry(tag(7)), tag(8)),
	}
	return tree
}

// levelCase is a filter at one level and how many of tree.lines it covers (lines[0:covers]).
type levelCase struct {
	name   string
	filter func(tree filterTree) dto.MachineFilter
	covers int
}

var filterLevelCases = []levelCase{
	{"no level (All machines of the company)", func(filterTree) dto.MachineFilter { return dto.MachineFilter{} }, 8},
	{"country", func(f filterTree) dto.MachineFilter { return dto.MachineFilter{UserCountryID: &f.country} }, 7},
	{"state", func(f filterTree) dto.MachineFilter {
		return dto.MachineFilter{UserCountryID: &f.country, UserStateID: &f.state}
	}, 6},
	{"city", func(f filterTree) dto.MachineFilter {
		return dto.MachineFilter{UserCountryID: &f.country, UserStateID: &f.state, UserCityID: &f.city}
	}, 5},
	{"branch", func(f filterTree) dto.MachineFilter {
		return dto.MachineFilter{UserCountryID: &f.country, UserStateID: &f.state, UserCityID: &f.city, BranchID: &f.branch}
	}, 4},
	{"location, all shops and lines (the screenshot)", func(f filterTree) dto.MachineFilter {
		return dto.MachineFilter{UserCountryID: &f.country, UserStateID: &f.state, UserCityID: &f.city, BranchID: &f.branch, UserLocationID: &f.location}
	}, 3},
	{"shop, all lines", func(f filterTree) dto.MachineFilter {
		return dto.MachineFilter{UserCountryID: &f.country, UserStateID: &f.state, UserCityID: &f.city, BranchID: &f.branch, UserLocationID: &f.location, UserShopID: &f.shop}
	}, 2},
	{"line", func(f filterTree) dto.MachineFilter {
		return dto.MachineFilter{UserCountryID: &f.country, UserStateID: &f.state, UserCityID: &f.city, BranchID: &f.branch, UserLocationID: &f.location, UserShopID: &f.shop, UserLineID: &f.line}
	}, 1},
}

// TestMachineFilterMatrix runs every combination of the user form's machine access through the
// real API. Everything it adds belongs to a throwaway customer and is soft-deleted at the end.
// Run with -v to see each result.
func TestMachineFilterMatrix(t *testing.T) {
	suffix := uniqueSuffix()
	company := newTestCompany(t, "X"+suffix, 0)
	other := newTestCompany(t, "Y"+suffix, 1)
	tree := newFilterTree(t, company, suffix)
	customer := newUserAccessRouter(actAs(company.customerID, model.RoleCustomer))

	addMachines := func(tag string) []int64 {
		ids := make([]int64, len(tree.lines))
		for index, lineID := range tree.lines {
			ids[index] = addTestMachine(t, company, fmt.Sprintf("MX-%s-%d-%s", tag, index, suffix), lineID)
		}
		return ids
	}
	existing := addMachines("old")

	body := userBody(t, company, "mx"+suffix, model.RoleSupervisor, nil)
	user := createUser(t, customer, body)
	userRouter := newUserAccessRouter(actAsMember(user.ID, model.RoleSupervisor, company.customerID))
	userPath := fmt.Sprintf("%s/users/%d", apiV1, user.ID)

	seen := func() []int64 { return listIDs(t, sendRequest(t, userRouter, http.MethodGet, apiV1+"/machines", nil)) }
	save := func(machineIDs []int64, filter *dto.MachineFilter) apiResponse {
		update := updateBody(body)
		update.MachineIDs = machineIDs
		update.MachineFilter = filter
		return sendRequest(t, customer, http.MethodPut, userPath, update)
	}
	mustSave := func(t *testing.T, machineIDs []int64, filter *dto.MachineFilter) {
		t.Helper()
		expectStatus(t, save(machineIDs, filter), http.StatusOK, "success")
	}
	withAll := func(filter dto.MachineFilter) *dto.MachineFilter {
		filter.AllMachines = true
		return &filter
	}
	sorted := func(ids ...[]int64) []int64 {
		merged := slices.Concat(ids...)
		slices.Sort(merged)
		return merged
	}

	// 1. All machines ticked: every filter level against a machine on every line
	for _, level := range filterLevelCases {
		t.Run("all machines ticked/"+level.name, func(t *testing.T) {
			mustSave(t, nil, withAll(level.filter(tree)))
			if got, want := seen(), sorted(existing[:level.covers]); !slices.Equal(got, want) {
				t.Errorf("sees %v; want %v", got, want)
			}
		})
	}

	// 2. All machines ticked: a machine added later appears without saving the user again
	later := addMachines("new")
	for _, level := range filterLevelCases {
		t.Run("machine added later/all ticked/"+level.name, func(t *testing.T) {
			mustSave(t, nil, withAll(level.filter(tree)))
			if got, want := seen(), sorted(existing[:level.covers], later[:level.covers]); !slices.Equal(got, want) {
				t.Errorf("sees %v; want %v", got, want)
			}
		})
	}

	// 3. Same filter, All machines NOT ticked, one machine picked (the screenshot): nothing is automatic
	for _, level := range filterLevelCases {
		t.Run("machine added later/not ticked, one picked/"+level.name, func(t *testing.T) {
			mustSave(t, []int64{existing[0]}, new(level.filter(tree)))
			got := seen()
			if want := []int64{existing[0]}; !slices.Equal(got, want) {
				t.Errorf("sees %v; want only the picked machine %v", got, want)
			}
			t.Logf("RESULT not ticked, filter %q: sees only the picked machine, %d machines under the filter are not given", level.name, level.covers*2-1)
		})
	}

	// 4. A new line, shop or location under the filter counts as well
	t.Run("new line under a shop filter is included", func(t *testing.T) {
		mustSave(t, nil, withAll(filterLevelCases[6].filter(tree)))
		line := insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, company.customerID, tree.shop, "NewLine"+suffix)
		machine := addTestMachine(t, company, "MX-newline-"+suffix, line)
		if got := seen(); !slices.Contains(got, machine) {
			t.Errorf("sees %v; want the machine %d on the new line", got, machine)
		}
	})
	t.Run("new shop and line under a location filter are included", func(t *testing.T) {
		mustSave(t, nil, withAll(filterLevelCases[5].filter(tree)))
		shop := insertTestRow(t, "user_shop", `INSERT INTO user_shop (user_id, user_location_id, name) VALUES ($1, $2, $3) RETURNING id`, company.customerID, tree.location, "NewShop"+suffix)
		line := insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, company.customerID, shop, "NewShopLine"+suffix)
		machine := addTestMachine(t, company, "MX-newshop-"+suffix, line)
		if got := seen(); !slices.Contains(got, machine) {
			t.Errorf("sees %v; want the machine %d under the new shop", got, machine)
		}
	})

	// 5. Different filter + different machines: can one user combine them?
	lineFilter := func(index int) dto.MachineFilter {
		return dto.MachineFilter{UserLineID: &tree.lines[index]}
	}
	t.Run("combining/one filter row only: a second All machines filter replaces the first", func(t *testing.T) {
		mustSave(t, nil, withAll(lineFilter(0)))
		first := seen()
		mustSave(t, nil, withAll(lineFilter(4)))
		second := seen()
		var rows int
		if err := testDBPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM user_machine_filters WHERE user_id = $1`, user.ID).Scan(&rows); err != nil {
			t.Fatal(err)
		}
		t.Logf("RESULT line A → %v, then line E → %v, filter rows = %d", first, second, rows)
		if rows != 1 || slices.Contains(second, existing[0]) {
			t.Errorf("want one filter row and line A's machines gone; rows %d, sees %v", rows, second)
		}
	})
	t.Run("combining/picks from several lines with no All: allowed", func(t *testing.T) {
		picks := []int64{existing[0], existing[4], existing[7]}
		mustSave(t, picks, &dto.MachineFilter{UserCountryID: &tree.country})
		got := seen()
		t.Logf("RESULT filter = country, picks on 3 lines (one in another country) → sees %v", got)
		if !slices.Equal(got, sorted(picks)) {
			t.Errorf("sees %v; want %v", got, sorted(picks))
		}
	})
	t.Run("combining/All on one line plus a pick from another line: the pick is dropped", func(t *testing.T) {
		result := save([]int64{existing[4]}, withAll(lineFilter(0)))
		got := seen()
		t.Logf("RESULT status %d %q, All on line A + pick on line E → sees %v, line E machine given: %v", result.StatusCode, result.Message, got, slices.Contains(got, existing[4]))
		if result.StatusCode != http.StatusOK {
			t.Errorf("status %d %q", result.StatusCode, result.Message)
		}
	})
	t.Run("combining/saving All removes the picks made before", func(t *testing.T) {
		mustSave(t, []int64{existing[4]}, &dto.MachineFilter{UserLineID: &tree.lines[4]})
		mustSave(t, nil, withAll(lineFilter(0)))
		var stored int
		if err := testDBPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM user_machine_access WHERE user_id = $1`, user.ID).Scan(&stored); err != nil {
			t.Fatal(err)
		}
		if got := seen(); stored != 0 || slices.Contains(got, existing[4]) {
			t.Errorf("stored picks %d, sees %v; want the old pick gone", stored, got)
		}
	})

	// 6. Machines that move, lose their line or are deleted
	t.Run("machine moved out of the filter loses access, moved in gains it", func(t *testing.T) {
		mustSave(t, nil, withAll(lineFilter(0)))
		moved := addTestMachine(t, company, "MX-moved-"+suffix, tree.lines[4])
		if slices.Contains(seen(), moved) {
			t.Fatalf("machine on another line is visible before the move")
		}
		setMachineLine(t, moved, tree.lines[0])
		if !slices.Contains(seen(), moved) {
			t.Errorf("machine moved onto the filtered line is not visible")
		}
		setMachineLine(t, moved, tree.lines[4])
		if slices.Contains(seen(), moved) {
			t.Errorf("machine moved off the filtered line is still visible")
		}
	})
	t.Run("picked machine moved out of the filter: keeps access", func(t *testing.T) {
		moved := addTestMachine(t, company, "MX-pickmoved-"+suffix, tree.lines[0])
		mustSave(t, []int64{moved}, &dto.MachineFilter{UserLineID: &tree.lines[0]})
		setMachineLine(t, moved, tree.lines[4])
		t.Logf("RESULT picked machine moved off the filter's line: still visible = %v", slices.Contains(seen(), moved))
		if !slices.Contains(seen(), moved) {
			t.Errorf("want a picked machine to stay visible")
		}
	})
	t.Run("machine without a line", func(t *testing.T) {
		result := sendRequest(t, newUserAccessRouter(asTestUser), http.MethodPost, apiV1+"/machines",
			with(machineBody(company.machineFixture, "MX-noline-"+suffix), "lineId", nil))
		expectStatus(t, result, http.StatusCreated, "created")
		var machine machineRecord
		decodeData(t, result, &machine)
		t.Cleanup(func() {
			_, _ = testDBPool.Exec(context.Background(), `UPDATE machines SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, machine.ID)
		})
		mustSave(t, nil, withAll(dto.MachineFilter{}))
		if !slices.Contains(seen(), machine.ID) {
			t.Errorf("All machines of the company should include a machine without a line")
		}
		mustSave(t, nil, withAll(dto.MachineFilter{UserCountryID: &tree.country}))
		if slices.Contains(seen(), machine.ID) {
			t.Errorf("a filter with a level should not include a machine without a line")
		}
	})
	t.Run("picked machine deleted: the user can still be saved", func(t *testing.T) {
		gone := addTestMachine(t, company, "MX-gone-"+suffix, tree.lines[0])
		keep := existing[0]
		mustSave(t, []int64{gone, keep}, &dto.MachineFilter{UserLineID: &tree.lines[0]})
		softDeleteMachine(t, gone)
		loaded := getUserAccess(t, customer, userPath)
		t.Logf("RESULT after deleting a picked machine, GET machineIds = %v", loaded.MachineIDs)
		if slices.Contains(loaded.MachineIDs, gone) {
			t.Errorf("deleted machine is still returned")
		}
		expectStatus(t, save(loaded.MachineIDs, &dto.MachineFilter{UserLineID: &tree.lines[0]}), http.StatusOK, "success")
	})

	// 7. A level in use cannot be deleted: the user would lose the machines and could not be saved
	organization := newFilterLevelRouter()
	deleteLevel := func(path string, id int64) apiResponse {
		return sendRequest(t, organization, http.MethodDelete, fmt.Sprintf("%s/%s/%d?userId=%d", apiV1, path, id, company.customerID), nil)
	}
	owner := company.customerID
	for _, level := range []struct {
		word, path string
		add        func(tag string) (int64, dto.MachineFilter)
	}{
		{"country", "user-country", func(tag string) (int64, dto.MachineFilter) {
			id := insertTestRow(t, "user_country", `INSERT INTO user_country (user_id, name) VALUES ($1, $2) RETURNING id`, owner, tag)
			return id, dto.MachineFilter{UserCountryID: &id}
		}},
		{"state", "user-state", func(tag string) (int64, dto.MachineFilter) {
			id := insertTestRow(t, "user_state", `INSERT INTO user_state (user_id, user_country_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, tree.country, tag)
			return id, dto.MachineFilter{UserStateID: &id}
		}},
		{"city", "user-city", func(tag string) (int64, dto.MachineFilter) {
			id := insertTestRow(t, "user_city", `INSERT INTO user_city (user_id, user_state_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, tree.state, tag)
			return id, dto.MachineFilter{UserCityID: &id}
		}},
		{"branch", "branches", func(tag string) (int64, dto.MachineFilter) {
			id := insertTestRow(t, "branches", `INSERT INTO branches (code, name, user_city_id) VALUES ($1, $2, $3) RETURNING id`, tag, tag, tree.city)
			return id, dto.MachineFilter{BranchID: &id}
		}},
		{"location", "user-location", func(tag string) (int64, dto.MachineFilter) {
			id := insertTestRow(t, "user_location", `INSERT INTO user_location (user_id, branch_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, tree.branch, tag)
			return id, dto.MachineFilter{UserLocationID: &id}
		}},
		{"shop", "user-shop", func(tag string) (int64, dto.MachineFilter) {
			id := insertTestRow(t, "user_shop", `INSERT INTO user_shop (user_id, user_location_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, tree.location, tag)
			return id, dto.MachineFilter{UserShopID: &id}
		}},
		{"line", "user-line", func(tag string) (int64, dto.MachineFilter) {
			id := insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, tree.shop, tag)
			return id, dto.MachineFilter{UserLineID: &id}
		}},
	} {
		t.Run("deleting a level in a filter is refused/"+level.word, func(t *testing.T) {
			want := "This " + level.word + " is used in a user's machine access. Change that user first"
			id, filter := level.add("Del" + level.word + suffix)
			mustSave(t, nil, withAll(filter))
			expectStatus(t, deleteLevel(level.path, id), http.StatusBadRequest, want)
			mustSave(t, nil, &filter) // All not ticked: re-saving the user would fail just the same
			expectStatus(t, deleteLevel(level.path, id), http.StatusBadRequest, want)
			mustSave(t, nil, nil)
			if result := deleteLevel(level.path, id); result.StatusCode != http.StatusOK {
				t.Errorf("after the filter moved on, delete → %d %q; want 200", result.StatusCode, result.Message)
			}
		})
	}
	t.Run("deleting a line with machines is refused", func(t *testing.T) {
		line := insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, tree.shop, "DelM"+suffix)
		machine := addTestMachine(t, company, "MX-delline-"+suffix, line)
		expectStatus(t, deleteLevel("user-line", line), http.StatusBadRequest, "This line has machines. Move or delete them first")
		setMachineLine(t, machine, tree.lines[0])
		if result := deleteLevel("user-line", line); result.StatusCode != http.StatusOK {
			t.Errorf("after the machine moved, delete → %d %q; want 200", result.StatusCode, result.Message)
		}
	})
	t.Run("a deleted machine does not keep its line", func(t *testing.T) {
		line := insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, tree.shop, "DelD"+suffix)
		softDeleteMachine(t, addTestMachine(t, company, "MX-delgone-"+suffix, line))
		if result := deleteLevel("user-line", line); result.StatusCode != http.StatusOK {
			t.Errorf("delete → %d %q; want 200", result.StatusCode, result.Message)
		}
	})
	t.Run("a deleted user's filter does not keep its line", func(t *testing.T) {
		line := insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, tree.shop, "DelU"+suffix)
		goneBody := userBody(t, company, "mxg"+suffix, model.RoleOperator, nil)
		goneBody.MachineFilter = withAll(dto.MachineFilter{UserLineID: &line})
		gone := createUser(t, customer, goneBody)
		expectStatus(t, deleteLevel("user-line", line), http.StatusBadRequest, "This line is used in a user's machine access. Change that user first")
		if _, err := testDBPool.Exec(context.Background(), `UPDATE users SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, gone.ID); err != nil {
			t.Fatal(err)
		}
		if result := deleteLevel("user-line", line); result.StatusCode != http.StatusOK {
			t.Errorf("delete → %d %q; want 200", result.StatusCode, result.Message)
		}
	})

	// 8. Bad input
	t.Run("bad input/another company's machine is refused", func(t *testing.T) {
		expectStatus(t, save([]int64{other.machineIDs[0]}, nil), http.StatusBadRequest, "Please select valid machines")
	})
	t.Run("bad input/another company's machine with All ticked", func(t *testing.T) {
		result := save([]int64{other.machineIDs[0]}, withAll(dto.MachineFilter{}))
		got := seen()
		t.Logf("RESULT All + another company's machine id → %d %q, given: %v", result.StatusCode, result.Message, slices.Contains(got, other.machineIDs[0]))
		if slices.Contains(got, other.machineIDs[0]) {
			t.Errorf("user got another company's machine")
		}
	})
	t.Run("bad input/another company's level is refused", func(t *testing.T) {
		expectStatus(t, save(nil, withAll(dto.MachineFilter{UserLineID: &other.lineID})), http.StatusBadRequest, "Please select valid machine filter")
	})
	t.Run("bad input/levels that are not one chain are refused", func(t *testing.T) {
		expectStatus(t, save(nil, withAll(dto.MachineFilter{UserShopID: &tree.shop, UserLineID: &tree.lines[4]})), http.StatusBadRequest, "Please select valid machine filter")
	})
	t.Run("bad input/unknown level is refused", func(t *testing.T) {
		missing := int64(2_000_000_000)
		expectStatus(t, save(nil, withAll(dto.MachineFilter{UserLineID: &missing})), http.StatusBadRequest, "Please select valid machine filter")
	})

	// 9. Other ways in
	t.Run("nothing picked and nothing ticked: no machines", func(t *testing.T) {
		mustSave(t, nil, nil)
		if got := seen(); len(got) != 0 {
			t.Errorf("sees %v; want none", got)
		}
	})
	t.Run("filter levels with nothing ticked or picked: no machines", func(t *testing.T) {
		mustSave(t, nil, new(lineFilter(0)))
		if got := seen(); len(got) != 0 {
			t.Errorf("sees %v; want none", got)
		}
	})
	t.Run("inactive machine is listed under All", func(t *testing.T) {
		mustSave(t, nil, withAll(lineFilter(0)))
		machine := addTestMachine(t, company, "MX-inactive-"+suffix, tree.lines[0])
		if _, err := testDBPool.Exec(context.Background(), `UPDATE machines SET status = 'inactive' WHERE id = $1`, machine); err != nil {
			t.Fatal(err)
		}
		t.Logf("RESULT inactive machine under the filter: visible = %v", slices.Contains(seen(), machine))
	})
	t.Run("role changed to admin clears the access, back to supervisor starts empty", func(t *testing.T) {
		mustSave(t, []int64{existing[0]}, &dto.MachineFilter{UserLineID: &tree.lines[0]})
		update := updateBody(body)
		update.RoleID = testRoleID(t, model.RoleAdmin)
		update.MachineIDs = []int64{existing[0]}
		update.MachineFilter = &dto.MachineFilter{UserLineID: &tree.lines[0]}
		expectStatus(t, sendRequest(t, customer, http.MethodPut, userPath, update), http.StatusOK, "success")
		loaded := getUserAccess(t, customer, userPath)
		if len(loaded.MachineIDs) != 0 || loaded.MachineFilter != nil {
			t.Errorf("admin keeps machineIds %v, filter %+v", loaded.MachineIDs, loaded.MachineFilter)
		}
		mustSave(t, nil, nil)
		if got := seen(); len(got) != 0 {
			t.Errorf("after going back sees %v; want none", got)
		}
	})
	t.Run("superadmin saves for a company", func(t *testing.T) {
		superadmin := newUserAccessRouter(asTestUser)
		superBody := userBody(t, company, "mxs"+suffix, model.RoleOperator, nil)
		superBody.CompanyID = &company.customerID
		superBody.MachineFilter = withAll(lineFilter(1))
		operator := createUser(t, superadmin, superBody)
		router := newUserAccessRouter(actAsMember(operator.ID, model.RoleOperator, company.customerID))
		got := listIDs(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines", nil))
		if want := sorted(existing[1:2], later[1:2]); !slices.Equal(got, want) {
			t.Errorf("operator sees %v; want %v", got, want)
		}
	})
	t.Run("one user's access does not touch another user's", func(t *testing.T) {
		mustSave(t, nil, withAll(lineFilter(0)))
		otherBody := userBody(t, company, "mxo"+suffix, model.RoleOperator, []int64{existing[4]})
		operator := createUser(t, customer, otherBody)
		router := newUserAccessRouter(actAsMember(operator.ID, model.RoleOperator, company.customerID))
		got := listIDs(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines", nil))
		if !slices.Equal(got, []int64{existing[4]}) {
			t.Errorf("operator sees %v; want only %d", got, existing[4])
		}
	})
}

func setMachineLine(t *testing.T, machineID, lineID int64) {
	t.Helper()
	if _, err := testDBPool.Exec(context.Background(), `UPDATE machines SET line_id = $2 WHERE id = $1`, machineID, lineID); err != nil {
		t.Fatalf("move machine: %v", err)
	}
}

func softDeleteMachine(t *testing.T, machineID int64) {
	t.Helper()
	if _, err := testDBPool.Exec(context.Background(), `UPDATE machines SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, machineID); err != nil {
		t.Fatalf("delete machine: %v", err)
	}
}

// newFilterLevelRouter is the real router with the organization levels' handlers, as the superadmin.
func newFilterLevelRouter() http.Handler {
	countries := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
	states := service.NewUserStateService(repository.NewUserStateRepository(testDBPool), countries)
	cities := service.NewUserCityService(repository.NewUserCityRepository(testDBPool), states)
	branches := service.NewBranchService(repository.NewBranchRepository(testDBPool), cities)
	locations := service.NewUserLocationService(repository.NewUserLocationRepository(testDBPool), branches)
	shops := service.NewUserShopService(repository.NewUserShopRepository(testDBPool), locations)
	lines := service.NewUserLineService(repository.NewUserLineRepository(testDBPool), shops)
	return routes.New(routes.Handlers{
		UserCountry:  handler.NewUserCountryHandler(countries),
		UserState:    handler.NewUserStateHandler(states),
		UserCity:     handler.NewUserCityHandler(cities),
		Branch:       handler.NewBranchHandler(branches),
		UserLocation: handler.NewUserLocationHandler(locations),
		UserShop:     handler.NewUserShopHandler(shops),
		UserLine:     handler.NewUserLineHandler(lines),
	}, asTestUser)
}

type userAccess struct {
	MachineIDs    []int64                  `json:"machineIds"`
	MachineFilter *model.UserMachineFilter `json:"machineFilter"`
}

func getUserAccess(t *testing.T, router http.Handler, path string) userAccess {
	t.Helper()
	result := sendRequest(t, router, http.MethodGet, path, nil)
	expectStatus(t, result, http.StatusOK, "success")
	var access userAccess
	decodeData(t, result, &access)
	return access
}
