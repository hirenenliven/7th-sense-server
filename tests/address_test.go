package tests

import (
	"context"
	"errors"
	"testing"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

func TestAddressCRUD(t *testing.T) {
	prepareProfileTable(t, "addresses")

	ctx := context.Background()
	addressService := service.NewAddressService(repository.NewAddressRepository(testDBPool))

	t.Run("get before create returns empty list", func(t *testing.T) {
		addresses, err := addressService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() error = %v", err)
		}
		if len(addresses) != 0 {
			t.Errorf("GetByUserID() = %d addresses; want 0", len(addresses))
		}
	})

	t.Run("create", func(t *testing.T) {
		created, err := addressService.Create(ctx, model.Address{
			UserID:        testUserID,
			AddressType:   "home",
			AddressLine1:  "123 Main Street",
			AddressLine2:  "Apt 4B",
			City:          "Rajkot",
			StateName:     "Gujarat",
			CountryName:   "India",
			PinCode:       "360001",
			ContactPerson: "John Doe",
			ContactPhone:  "+919876543210",
			IsDefault:     true,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if created.ID == 0 || created.UserID != testUserID {
			t.Errorf("Create() id = %d, userId = %d; want non-zero id, userId %d", created.ID, created.UserID, testUserID)
		}
		if created.AddressType != "home" || created.AddressLine1 != "123 Main Street" ||
			created.City != "Rajkot" || created.PinCode != "360001" {
			t.Errorf("Create() = %+v; fields not saved", created)
		}
		if created.IsDeleted {
			t.Error("Create() isDeleted = true; want false")
		}
	})

	t.Run("read", func(t *testing.T) {
		addresses, err := addressService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() error = %v", err)
		}
		if len(addresses) != 1 {
			t.Fatalf("GetByUserID() = %d addresses; want 1", len(addresses))
		}
		if addresses[0].AddressLine1 != "123 Main Street" || addresses[0].City != "Rajkot" {
			t.Errorf("GetByUserID() = %+v; want the created address", addresses[0])
		}
	})

	t.Run("update", func(t *testing.T) {
		addresses, err := addressService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() error = %v", err)
		}
		if len(addresses) == 0 {
			t.Fatal("GetByUserID() returned no addresses")
		}

		updated, err := addressService.Update(ctx, model.Address{
			ID:            addresses[0].ID,
			UserID:        testUserID,
			AddressType:   "work",
			AddressLine1:  "456 Business Park",
			AddressLine2:  "", // cleared
			City:          "Ahmedabad",
			StateName:     "Gujarat",
			CountryName:   "India",
			PinCode:       "380001",
			ContactPerson: "Jane Smith",
			ContactPhone:  "+919876543211",
			IsDefault:     false,
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		if updated.AddressType != "work" || updated.AddressLine1 != "456 Business Park" ||
			updated.City != "Ahmedabad" || updated.PinCode != "380001" {
			t.Errorf("Update() = %+v; fields not updated", updated)
		}

		addressList, err := addressService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() after update error = %v", err)
		}
		if addressList[0].AddressLine1 != "456 Business Park" {
			t.Errorf("GetByUserID() addressLine1 = %q; want updated value", addressList[0].AddressLine1)
		}
	})

	t.Run("delete", func(t *testing.T) {
		addresses, err := addressService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() error = %v", err)
		}
		if len(addresses) == 0 {
			t.Fatal("GetByUserID() returned no addresses")
		}

		if err := addressService.Delete(ctx, addresses[0].ID); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		addressList, err := addressService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() after delete error = %v", err)
		}
		if len(addressList) != 0 {
			t.Errorf("GetByUserID() after delete = %d addresses; want 0", len(addressList))
		}
		if err := addressService.Delete(ctx, addresses[0].ID); !errors.Is(err, response.NotFound) {
			t.Errorf("second Delete() error = %v; want %v", err, response.NotFound)
		}
	})
}
