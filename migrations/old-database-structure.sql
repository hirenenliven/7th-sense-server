-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Oct 30, 2026 at 09:54 AM
-- Server version: 8.0.30
-- PHP Version: 8.1.10

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `machine_monitor`
--

-- --------------------------------------------------------

--
-- Table structure for table `batch_tbl`
--

CREATE TABLE `batch_tbl` (
  `batchId` int NOT NULL,
  `machineId` int NOT NULL DEFAULT '0',
  `cycleId` varchar(255) NOT NULL,
  `batchName` varchar(255) NOT NULL,
  `createdDate` datetime NOT NULL,
  `updatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `branches`
--

CREATE TABLE `branches` (
  `iBranchId` int NOT NULL,
  `vBranchCode` varchar(100) NOT NULL,
  `iBranchGroup` int NOT NULL,
  `iRegionId` int NOT NULL,
  `iParentBranch` int NOT NULL,
  `vBranchName` varchar(100) NOT NULL,
  `vBranchEmail` varchar(100) NOT NULL,
  `vColorCode` varchar(10) NOT NULL,
  `eGstDetails` tinyint NOT NULL,
  `vGstBranchName` varchar(150) NOT NULL,
  `vAddress` varchar(250) NOT NULL,
  `iStateId` int NOT NULL,
  `vCity` varchar(150) NOT NULL,
  `vStateCode` varchar(15) NOT NULL,
  `vGstIn` varchar(30) NOT NULL,
  `vPanNumber` varchar(30) NOT NULL,
  `eStatus` enum('y','n') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `brand`
--

CREATE TABLE `brand` (
  `id` int NOT NULL,
  `name` varchar(100) NOT NULL,
  `del_status` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `breaks`
--

CREATE TABLE `breaks` (
  `iBreakId` int NOT NULL,
  `iUserId` int NOT NULL,
  `vBreakName` varchar(100) NOT NULL,
  `tStartTime` time NOT NULL,
  `tEndTime` time NOT NULL,
  `dTotalHours` double(5,3) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `category`
--

CREATE TABLE `category` (
  `id` int NOT NULL,
  `iMCategoryId` int NOT NULL,
  `name` varchar(100) NOT NULL,
  `iConfigIds` varchar(100) NOT NULL,
  `del_status` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ci_sessions`
--

CREATE TABLE `ci_sessions` (
  `id` varchar(128) NOT NULL,
  `ip_address` varchar(45) NOT NULL,
  `timestamp` int UNSIGNED NOT NULL DEFAULT '0',
  `data` blob NOT NULL,
  `iUserId` int DEFAULT NULL,
  `dCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `clientele`
--

CREATE TABLE `clientele` (
  `iClienteleId` int NOT NULL,
  `vClientNumber` varchar(250) NOT NULL,
  `vClienteleName` varchar(250) NOT NULL,
  `vEmail` varchar(250) NOT NULL,
  `vContactNumber` varchar(250) NOT NULL,
  `vAddress` varchar(250) NOT NULL,
  `iUserId` int NOT NULL,
  `iCountryId` int NOT NULL,
  `iStateId` int NOT NULL,
  `vCityName` varchar(250) NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `controller`
--

CREATE TABLE `controller` (
  `id` int NOT NULL,
  `iProductIds` varchar(250) NOT NULL,
  `controller_variant` varchar(100) NOT NULL,
  `brand_id` int NOT NULL,
  `del_status` int NOT NULL,
  `tDependents` text NOT NULL,
  `tDependentsOpt` longtext NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `downtime_items`
--

CREATE TABLE `downtime_items` (
  `iItemId` int NOT NULL,
  `iUserId` int NOT NULL,
  `dDowntimeDate` date NOT NULL,
  `iMachineId` int NOT NULL,
  `iShiftId` int DEFAULT NULL,
  `tDownStart` time NOT NULL,
  `tDownEnd` time NOT NULL,
  `dTotalHours` float NOT NULL,
  `vReason` varchar(200) NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL,
  `vOperator` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `downtime_reason`
--

CREATE TABLE `downtime_reason` (
  `iDtrId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iDataId` int NOT NULL,
  `iDownReason` int NOT NULL,
  `iItemId` int NOT NULL,
  `dDownTimeDate` date NOT NULL,
  `dDTStartTime` time NOT NULL,
  `dDTEndTime` time NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL DEFAULT 'y'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `dropdown_key`
--

CREATE TABLE `dropdown_key` (
  `iKeyId` int NOT NULL,
  `vKeyName` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `dropdown_value`
--

CREATE TABLE `dropdown_value` (
  `iValueId` int NOT NULL,
  `vValueName` varchar(100) NOT NULL,
  `iKeyId` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `helpenter_pages`
--

CREATE TABLE `helpenter_pages` (
  `iHcId` int NOT NULL,
  `vTitle` varchar(200) NOT NULL,
  `vSlug` varchar(200) NOT NULL,
  `tDescription` text NOT NULL,
  `iParent` int NOT NULL,
  `iSubParent` int NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL DEFAULT 'y',
  `dCreatedDate` datetime NOT NULL,
  `vLang` varchar(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `helpenter_pages_bkup`
--

CREATE TABLE `helpenter_pages_bkup` (
  `iHcId` int NOT NULL,
  `vTitle` varchar(200) NOT NULL,
  `vSlug` varchar(200) NOT NULL,
  `tDescription` text NOT NULL,
  `iParent` int NOT NULL,
  `iSubParent` int NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL DEFAULT 'y',
  `dCreatedDate` datetime NOT NULL,
  `vLang` varchar(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `idle_reason`
--

CREATE TABLE `idle_reason` (
  `iIdlerId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iDataId` int NOT NULL,
  `iIdleReason` int NOT NULL,
  `iItemId` int NOT NULL,
  `dDownTimeDate` date NOT NULL,
  `dDTStartTime` time NOT NULL,
  `dDTEndTime` time NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machines_data`
--

CREATE TABLE `machines_data` (
  `iDataId` int NOT NULL,
  `iModeId` int NOT NULL,
  `iSubModeId` longtext NOT NULL,
  `iMachineId` int NOT NULL,
  `iSeconds` int NOT NULL DEFAULT '1',
  `dCreatedDate` datetime NOT NULL,
  `lastUpDate` datetime NOT NULL,
  `isLoadingOver` tinyint(1) NOT NULL,
  `Remark` text NOT NULL,
  `dActCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `isUpdateFinix` tinyint NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machines_data_bkp`
--

CREATE TABLE `machines_data_bkp` (
  `iDataId` int NOT NULL,
  `iModeId` int NOT NULL,
  `iSubModeId` longtext NOT NULL,
  `iMachineId` int NOT NULL,
  `iSeconds` int NOT NULL DEFAULT '1',
  `dCreatedDate` datetime NOT NULL,
  `lastUpDate` datetime NOT NULL,
  `isLoadingOver` tinyint(1) NOT NULL,
  `Remark` text NOT NULL,
  `dActCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machines_data_history`
--

CREATE TABLE `machines_data_history` (
  `iDataId` int NOT NULL,
  `iModeId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iSeconds` int NOT NULL DEFAULT '1',
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machines_data_temp`
--

CREATE TABLE `machines_data_temp` (
  `iDataId` int NOT NULL,
  `iModeId` int NOT NULL,
  `iSubModeId` longtext NOT NULL,
  `iMachineId` int NOT NULL,
  `iSeconds` int NOT NULL DEFAULT '1',
  `dCreatedDate` datetime NOT NULL,
  `lastUpDate` datetime NOT NULL,
  `isLoadingOver` tinyint(1) NOT NULL,
  `Remark` text NOT NULL,
  `dActCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `isUpdateFinix` tinyint NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machines_modes`
--

CREATE TABLE `machines_modes` (
  `iModeId` int NOT NULL,
  `vModeName` varchar(100) NOT NULL,
  `eModeFor` enum('Machine','User') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_accuracy`
--

CREATE TABLE `machine_accuracy` (
  `iAccuracyId` int NOT NULL,
  `vAccuracy` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_alarms`
--

CREATE TABLE `machine_alarms` (
  `iAlarmId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iAlarmCode` int NOT NULL,
  `vAlarmText` varchar(250) NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` tinyint(1) NOT NULL DEFAULT '1',
  `iHAl1` int NOT NULL,
  `iHAl5` int NOT NULL,
  `iHAl6` int NOT NULL,
  `iHAl7` int NOT NULL,
  `iHAl8` int NOT NULL,
  `iHAl9` int NOT NULL,
  `iHAl10` int NOT NULL,
  `iHAl11` int NOT NULL,
  `iHAl12` int NOT NULL,
  `iHAl13` int NOT NULL,
  `iHAl14` int NOT NULL,
  `iHAl15` int NOT NULL,
  `iHAl16` int NOT NULL,
  `iHAl17` int NOT NULL,
  `iHAl18` int NOT NULL,
  `iHAl19` int NOT NULL,
  `iHAl20` int NOT NULL,
  `iHAl21` int NOT NULL,
  `iHAl22` int NOT NULL,
  `iHAl23` int NOT NULL,
  `SpindleLoad` varchar(25) NOT NULL,
  `SpindleRpm` varchar(25) NOT NULL,
  `speedOvr` varchar(25) NOT NULL,
  `feedRateOvr` varchar(25) NOT NULL,
  `FeedRev` varchar(25) NOT NULL,
  `actFeedRate` varchar(25) NOT NULL,
  `axPowerX` varchar(25) NOT NULL,
  `axPowerY` varchar(25) NOT NULL,
  `axPowerZ` varchar(25) NOT NULL,
  `axLoadZ` varchar(25) NOT NULL,
  `axLoadY` varchar(25) NOT NULL,
  `axLoadX` varchar(25) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_alarms_history`
--

CREATE TABLE `machine_alarms_history` (
  `iAlarmId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iAlarmCode` int NOT NULL,
  `vAlarmText` varchar(250) NOT NULL,
  `isHealth` tinyint(1) NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_data_count`
--

CREATE TABLE `machine_data_count` (
  `iMdcId` int NOT NULL,
  `iModeId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iTotal` int NOT NULL,
  `iItemId` int NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_data_count_history`
--

CREATE TABLE `machine_data_count_history` (
  `iMdcId` int NOT NULL,
  `iModeId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iTotal` int NOT NULL,
  `iItemId` int NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_data_count_hourly`
--

CREATE TABLE `machine_data_count_hourly` (
  `iMdcId` int NOT NULL,
  `iModeId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iTotal` int NOT NULL,
  `iItemId` int NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_data_count_hourly_history`
--

CREATE TABLE `machine_data_count_hourly_history` (
  `iMdcId` int NOT NULL,
  `iModeId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iTotal` int NOT NULL,
  `iItemId` int NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_lineno`
--

CREATE TABLE `machine_lineno` (
  `iLineId` int NOT NULL,
  `vLineNo` varchar(10) NOT NULL,
  `iUserId` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_online_log`
--

CREATE TABLE `machine_online_log` (
  `iOnLogid` int NOT NULL,
  `iMachineId` int NOT NULL,
  `dOnlineDate` datetime NOT NULL,
  `dOffileDate` datetime NOT NULL,
  `eStatus` enum('0','1') NOT NULL COMMENT '1=ONLINE,0=OFFLINE'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `machine_programs`
--

CREATE TABLE `machine_programs` (
  `iProgramId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `vProgramName` varchar(150) NOT NULL,
  `vMainProgramName` varchar(150) NOT NULL,
  `vActiveTool` varchar(50) NOT NULL,
  `iopMode` tinyint(1) NOT NULL COMMENT '0 = JOG 1 = MDI 2 = AUTO',
  `dCreatedDate` datetime NOT NULL,
  `dToolUpdate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `maintenance_tasks`
--

CREATE TABLE `maintenance_tasks` (
  `iMTaskId` int NOT NULL,
  `vDescription` varchar(100) NOT NULL,
  `eStatus` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_api_timelog`
--

CREATE TABLE `mm_api_timelog` (
  `Id` int NOT NULL,
  `vTimeTaken` longtext NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_cycletime`
--

CREATE TABLE `mm_cycletime` (
  `iDataId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iModeId` int NOT NULL,
  `vProgramName` varchar(255) NOT NULL,
  `iItemId` int NOT NULL,
  `iSeconds` int NOT NULL DEFAULT '1',
  `iPartCount` int NOT NULL DEFAULT '0',
  `iInProduction` int NOT NULL,
  `dProdStartDate` datetime DEFAULT NULL,
  `iInHold` int NOT NULL,
  `dHoldStartDate` datetime DEFAULT NULL,
  `iInRapid` int NOT NULL,
  `dRapidStartDate` datetime DEFAULT NULL,
  `iInLdUld` int NOT NULL,
  `dLdUldStartDate` datetime DEFAULT NULL,
  `iCstatus` tinyint(1) NOT NULL DEFAULT '0' COMMENT ' 0-Cycle Runig,1- Unfinished, 2 - Cycle ended',
  `iCycleId` bigint NOT NULL,
  `dModUpDate` datetime NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `dCendTime` datetime NOT NULL,
  `Remark` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_cycletime_history`
--

CREATE TABLE `mm_cycletime_history` (
  `iDataId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `vProgramName` varchar(255) NOT NULL,
  `iItemId` int NOT NULL,
  `iPartCount` int NOT NULL DEFAULT '0',
  `iInProduction` int NOT NULL,
  `iInHold` int NOT NULL,
  `iInRapid` int NOT NULL,
  `iInLdUld` int NOT NULL,
  `iCstatus` tinyint(1) NOT NULL DEFAULT '0' COMMENT ' 0-Cycle Runig,1- Unfinished, 2 - Cycle ended',
  `iCycleId` bigint NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `dCendTime` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_daily_production`
--

CREATE TABLE `mm_daily_production` (
  `iDlyPrdId` int NOT NULL,
  `iItemId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iQty` int NOT NULL,
  `iGoodParts` int NOT NULL,
  `iScrapParts` int NOT NULL,
  `iSApproval` tinyint(1) NOT NULL COMMENT '1 = Approved, 0 - Not Approved',
  `iQApproval` tinyint(1) NOT NULL COMMENT '1 = Approved, 0 - Not Approved',
  `eStatus` tinyint(1) NOT NULL COMMENT '0 - None, 1 - Production Started, 2 Production Stoped, 3 - Production Paushed',
  `dPstartTime` datetime NOT NULL,
  `dPstopTime` datetime NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_daily_production_history`
--

CREATE TABLE `mm_daily_production_history` (
  `iDphId` int NOT NULL,
  `iItemId` int NOT NULL,
  `iGoodParts` int NOT NULL,
  `iScrapParts` int NOT NULL,
  `iUserId` int NOT NULL,
  `eActionType` enum('Supervisor','Quality','Operator') CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT 'Operator',
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_health`
--

CREATE TABLE `mm_health` (
  `iHealthId` int NOT NULL,
  `vHealthSenser` varchar(100) NOT NULL,
  `eType` tinyint(1) NOT NULL COMMENT '1=Alarm'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_log`
--

CREATE TABLE `mm_log` (
  `iLogId` bigint NOT NULL,
  `iMachineId` int NOT NULL,
  `incval` int NOT NULL,
  `vVariableKey` varchar(25) NOT NULL,
  `iVarId` int NOT NULL,
  `aBothValue` varchar(250) NOT NULL,
  `dCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `eDtaken` enum('Yes','No') NOT NULL DEFAULT 'No'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_log_fanuc`
--

CREATE TABLE `mm_log_fanuc` (
  `iLogId` bigint NOT NULL,
  `iMachineId` int NOT NULL,
  `vVariableKey` varchar(25) NOT NULL,
  `vValue` varchar(250) NOT NULL,
  `dCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `eDtaken` enum('Yes','No') NOT NULL DEFAULT 'No'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_machines`
--

CREATE TABLE `mm_machines` (
  `iMachineId` int NOT NULL,
  `iUserId` int NOT NULL,
  `iMachineType` int NOT NULL,
  `iProductId` int NOT NULL,
  `vMachineName` varchar(50) NOT NULL,
  `vMachineSrNo` varchar(50) NOT NULL,
  `vLiveIp` varchar(20) NOT NULL,
  `iPort` int NOT NULL,
  `vUserName` varchar(15) NOT NULL,
  `vPaSSCode` varchar(25) NOT NULL,
  `vLocalIp` varchar(20) NOT NULL,
  `iBrandId` int NOT NULL,
  `iControllerId` int NOT NULL,
  `vControllerVer` varchar(10) NOT NULL,
  `iAccuracyId` int NOT NULL,
  `iLineId` int NOT NULL,
  `dInstDate` date NOT NULL,
  `dBillDate` date NOT NULL,
  `dWarrentyDate` date NOT NULL,
  `vAmc` varchar(50) NOT NULL,
  `eInsured` enum('y','n') NOT NULL,
  `eOnlineStatus` tinyint(1) NOT NULL COMMENT '1=ONLINE,0=OFFLINE',
  `eMaitananceReq` enum('y','n') NOT NULL,
  `dHourlyCost` double(10,2) NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL,
  `eIsMonitor` enum('y','n') NOT NULL,
  `dOnlineDate` date NOT NULL,
  `dOfflineDate` datetime NOT NULL,
  `iModeId` int NOT NULL,
  `eControllVer` tinyint(1) NOT NULL COMMENT '1=828D,2=840D',
  `aHealthIds` varchar(100) NOT NULL,
  `vLongitude` varchar(50) NOT NULL,
  `vLatitude` varchar(50) NOT NULL,
  `eCycleAnalysis` enum('Yes','No') NOT NULL DEFAULT 'No',
  `eMachineSystemType` enum('NewPlc','OldPlc') CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL DEFAULT 'NewPlc',
  `eSecurityMode` enum('None','Basic256Sha256') NOT NULL DEFAULT 'Basic256Sha256'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_machines_bkp`
--

CREATE TABLE `mm_machines_bkp` (
  `iMachineId` int NOT NULL,
  `iUserId` int NOT NULL,
  `iMachineType` int NOT NULL,
  `iProductId` int NOT NULL,
  `vMachineName` varchar(50) NOT NULL,
  `vMachineSrNo` varchar(50) NOT NULL,
  `vLiveIp` varchar(20) NOT NULL,
  `iPort` int NOT NULL,
  `vUserName` varchar(15) NOT NULL,
  `vPaSSCode` varchar(25) NOT NULL,
  `vLocalIp` varchar(20) NOT NULL,
  `iBrandId` int NOT NULL,
  `iControllerId` int NOT NULL,
  `vControllerVer` varchar(10) NOT NULL,
  `iAccuracyId` int NOT NULL,
  `iLineId` int NOT NULL,
  `dInstDate` date NOT NULL,
  `dBillDate` date NOT NULL,
  `dWarrentyDate` date NOT NULL,
  `vAmc` varchar(50) NOT NULL,
  `eInsured` enum('y','n') NOT NULL,
  `eOnlineStatus` tinyint(1) NOT NULL COMMENT '1=ONLINE,0=OFFLINE',
  `eMaitananceReq` enum('y','n') NOT NULL,
  `dHourlyCost` double(10,2) NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL,
  `eIsMonitor` enum('y','n') NOT NULL,
  `dOnlineDate` date NOT NULL,
  `dOfflineDate` datetime NOT NULL,
  `iModeId` int NOT NULL,
  `eControllVer` tinyint(1) NOT NULL COMMENT '1=828D,2=840D',
  `aHealthIds` varchar(100) NOT NULL,
  `vLongitude` varchar(50) NOT NULL,
  `vLatitude` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_machines_ori`
--

CREATE TABLE `mm_machines_ori` (
  `iMachineId` int NOT NULL,
  `iUserId` int NOT NULL,
  `iMachineType` int NOT NULL,
  `iProductId` int NOT NULL,
  `vMachineName` varchar(50) NOT NULL,
  `vMachineSrNo` varchar(50) NOT NULL,
  `vLiveIp` varchar(20) NOT NULL,
  `iPort` int NOT NULL,
  `vUserName` varchar(15) NOT NULL,
  `vPaSSCode` varchar(25) NOT NULL,
  `vLocalIp` varchar(20) NOT NULL,
  `iBrandId` int NOT NULL,
  `iControllerId` int NOT NULL,
  `vControllerVer` varchar(10) NOT NULL,
  `iAccuracyId` int NOT NULL,
  `iLineId` int NOT NULL,
  `dInstDate` date NOT NULL,
  `dBillDate` date NOT NULL,
  `dWarrentyDate` date NOT NULL,
  `vAmc` varchar(50) NOT NULL,
  `eInsured` enum('y','n') NOT NULL,
  `eOnlineStatus` tinyint(1) NOT NULL COMMENT '1=ONLINE,0=OFFLINE',
  `eMaitananceReq` enum('y','n') NOT NULL,
  `dHourlyCost` double(10,2) NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL,
  `eIsMonitor` enum('y','n') NOT NULL,
  `dOnlineDate` date NOT NULL,
  `dOfflineDate` datetime NOT NULL,
  `iModeId` int NOT NULL,
  `eControllVer` tinyint(1) NOT NULL COMMENT '1=828D,2=840D',
  `aHealthIds` varchar(100) NOT NULL,
  `vLongitude` varchar(50) NOT NULL,
  `vLatitude` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_modules`
--

CREATE TABLE `mm_modules` (
  `iModuleId` int NOT NULL,
  `vModuleName` varchar(300) DEFAULT NULL,
  `vConrollerName` varchar(165) DEFAULT NULL,
  `eStatus` tinyint(1) DEFAULT NULL,
  `iPermInputs` int DEFAULT NULL,
  `iPriority` int DEFAULT NULL,
  `isChart` tinyint(1) DEFAULT NULL,
  `isCritical` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_modules_old`
--

CREATE TABLE `mm_modules_old` (
  `iModuleId` int NOT NULL,
  `vModuleName` varchar(100) NOT NULL,
  `vConrollerName` varchar(55) NOT NULL,
  `eStatus` tinyint(1) NOT NULL DEFAULT '1',
  `iPermInputs` int NOT NULL,
  `iPriority` int NOT NULL,
  `isChart` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 - Not  a Chart, 1 - its chart page'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_operator_tasks`
--

CREATE TABLE `mm_operator_tasks` (
  `iOpTaskId` int NOT NULL,
  `iMachineType` varchar(255) NOT NULL,
  `iMTaskId` int NOT NULL,
  `iUserId` int DEFAULT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_programs`
--

CREATE TABLE `mm_programs` (
  `iPrgId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `vProgramName` varchar(200) NOT NULL,
  `dStartDate` datetime NOT NULL,
  `dEndDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_programs_history`
--

CREATE TABLE `mm_programs_history` (
  `iPrgId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `vProgramName` varchar(200) NOT NULL,
  `dStartDate` datetime NOT NULL,
  `dEndDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_spindleload`
--

CREATE TABLE `mm_spindleload` (
  `iTempId` varchar(25) NOT NULL,
  `iMachineId` int NOT NULL,
  `eTempType` int NOT NULL,
  `iLoad` float NOT NULL,
  `dCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` int DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_temperature`
--

CREATE TABLE `mm_temperature` (
  `iMachineId` int NOT NULL,
  `iTempId` bigint NOT NULL,
  `eTempType` int NOT NULL COMMENT '29=spindleBearingTemp,30=spidleVibration,31=HeadstockTemp',
  `iTemp` varchar(40) NOT NULL,
  `dCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_tool_life`
--

CREATE TABLE `mm_tool_life` (
  `iTlifeId` int NOT NULL,
  `vToolNo` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `iMachineId` int NOT NULL,
  `iItemId` int NOT NULL,
  `iCycleId` bigint NOT NULL,
  `iToollife` int NOT NULL,
  `dCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_tool_life_history`
--

CREATE TABLE `mm_tool_life_history` (
  `iTlifeId` int NOT NULL,
  `vToolNo` varchar(100) NOT NULL,
  `iMachineId` int NOT NULL,
  `iItemId` int NOT NULL,
  `iCycleId` bigint NOT NULL,
  `iToollife` int NOT NULL,
  `dCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mm_variables`
--

CREATE TABLE `mm_variables` (
  `iVarId` int NOT NULL,
  `vVarName` varchar(50) NOT NULL,
  `vVarPathv1` varchar(150) NOT NULL COMMENT '828D',
  `vVarPathv2` varchar(150) NOT NULL COMMENT '840D',
  `vRemarks` varchar(250) NOT NULL,
  `eStatus` enum('y','n') NOT NULL,
  `isStore` tinyint(1) NOT NULL,
  `iInterval` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='JYOTI MACHINE SHOP VARIABLE';

-- --------------------------------------------------------

--
-- Table structure for table `mm_variables_01-11-2022`
--

CREATE TABLE `mm_variables_01-11-2022` (
  `iVarId` int NOT NULL,
  `vVarName` varchar(50) NOT NULL,
  `vVarPathv1` varchar(150) NOT NULL COMMENT '828D',
  `vVarPathv2` varchar(150) NOT NULL COMMENT '840D',
  `vRemarks` varchar(250) NOT NULL,
  `eStatus` enum('y','n') NOT NULL,
  `isStore` tinyint(1) NOT NULL,
  `iInterval` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='JYOTI MACHINE SHOP VARIABLE';

-- --------------------------------------------------------

--
-- Table structure for table `mm_variables_21-09-2022`
--

CREATE TABLE `mm_variables_21-09-2022` (
  `iVarId` int NOT NULL,
  `vVarName` varchar(50) NOT NULL,
  `vVarPathv1` varchar(150) NOT NULL COMMENT '828D',
  `vVarPathv2` varchar(150) NOT NULL COMMENT '840D',
  `vRemarks` varchar(250) NOT NULL,
  `eStatus` enum('y','n') NOT NULL,
  `isStore` tinyint(1) NOT NULL,
  `iInterval` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='JYOTI MACHINE SHOP VARIABLE';

-- --------------------------------------------------------

--
-- Table structure for table `mm_variables_final`
--

CREATE TABLE `mm_variables_final` (
  `iVarId` int NOT NULL,
  `vVarName` varchar(50) NOT NULL,
  `vVarPathv1` varchar(150) NOT NULL COMMENT '828D',
  `vVarPathv2` varchar(150) NOT NULL COMMENT '840D',
  `vRemarks` varchar(250) NOT NULL,
  `eStatus` enum('y','n') NOT NULL,
  `isStore` tinyint(1) NOT NULL,
  `iInterval` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='FINAL LIST OF VARIABLE';

-- --------------------------------------------------------

--
-- Table structure for table `module_permission`
--

CREATE TABLE `module_permission` (
  `id` int NOT NULL,
  `iUserId` int NOT NULL,
  `iModuleId` int NOT NULL,
  `iPermission` tinyint(1) NOT NULL COMMENT '1=>read,Write=>2,Delete=>3	'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mst_countries`
--

CREATE TABLE `mst_countries` (
  `iCountryId` int NOT NULL,
  `vSortName` varchar(3) NOT NULL,
  `vCountryName` varchar(100) NOT NULL,
  `vPhoneCode` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- --------------------------------------------------------

--
-- Table structure for table `mst_downtime_reason`
--

CREATE TABLE `mst_downtime_reason` (
  `iDownReason` int NOT NULL,
  `iUserId` int NOT NULL,
  `vReasonCode` varchar(10) NOT NULL,
  `vReason` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mst_downtime_reason_copy`
--

CREATE TABLE `mst_downtime_reason_copy` (
  `iDownReason` int NOT NULL,
  `iUserId` int NOT NULL,
  `vReasonCode` varchar(10) NOT NULL,
  `vReason` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mst_idle_reason`
--

CREATE TABLE `mst_idle_reason` (
  `iIdleReason` int NOT NULL,
  `iUserId` int NOT NULL,
  `vReasonCode` varchar(10) NOT NULL,
  `vReason` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `mst_imagethumb`
--

CREATE TABLE `mst_imagethumb` (
  `iThumbId` int UNSIGNED NOT NULL,
  `vFolder` varchar(50) NOT NULL,
  `iHeight` int UNSIGNED NOT NULL,
  `iWidth` int UNSIGNED NOT NULL,
  `vDefaultImage` varchar(100) NOT NULL,
  `vIsactive` enum('y','n') NOT NULL DEFAULT 'y'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- --------------------------------------------------------

--
-- Table structure for table `mst_site_settings`
--

CREATE TABLE `mst_site_settings` (
  `iFieldId` int UNSIGNED NOT NULL,
  `vLabel` varchar(100) DEFAULT NULL,
  `vType` varchar(25) DEFAULT NULL,
  `vConstant` varchar(50) DEFAULT NULL,
  `vOptions` text NOT NULL,
  `vClass` varchar(100) DEFAULT NULL,
  `eRequired` enum('1','0') NOT NULL DEFAULT '0',
  `vValue` text NOT NULL,
  `vHint` varchar(255) DEFAULT NULL,
  `eEditable` enum('y','n') NOT NULL DEFAULT 'y',
  `eStatus` enum('y','n') NOT NULL DEFAULT 'n',
  `dUpdatedDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- --------------------------------------------------------

--
-- Table structure for table `mst_states`
--

CREATE TABLE `mst_states` (
  `iStateId` int NOT NULL,
  `StateName` varchar(30) NOT NULL,
  `iCountryId` int NOT NULL DEFAULT '1',
  `iStateCode` tinyint NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `oper_mainten_tasks`
--

CREATE TABLE `oper_mainten_tasks` (
  `iOpMtskId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iMTaskId` int NOT NULL,
  `eStatus` tinyint(1) NOT NULL COMMENT '0 - NOT Checked 1 - Checked',
  `dCheckedDate` datetime NOT NULL,
  `iItemId` int NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `parts`
--

CREATE TABLE `parts` (
  `iPartId` int NOT NULL,
  `iUserId` int NOT NULL,
  `iClienteleId` int NOT NULL,
  `vPartNo` varchar(100) NOT NULL,
  `vPartName` varchar(200) NOT NULL,
  `iPartCount` int NOT NULL DEFAULT '1',
  `tProductionTime` int NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `parts_process`
--

CREATE TABLE `parts_process` (
  `iProcessId` int NOT NULL,
  `iPartId` int NOT NULL,
  `vProcess` varchar(100) NOT NULL,
  `vProgramName` varchar(100) NOT NULL,
  `vApplicableMachine` varchar(100) NOT NULL,
  `vDrawingNo` varchar(100) NOT NULL,
  `vRevisionNo` varchar(100) NOT NULL,
  `vDrawingPdf` text NOT NULL,
  `iAccuracyId` varchar(10) NOT NULL,
  `iSetupTime` int NOT NULL,
  `iCycleTime` int NOT NULL,
  `iLUTime` int NOT NULL,
  `tProductionTime` int NOT NULL,
  `eSubProgram` enum('yes','no') NOT NULL DEFAULT 'no',
  `dCreatedDate` datetime NOT NULL,
  `vSubPartName` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `parts_sub_process`
--

CREATE TABLE `parts_sub_process` (
  `iPrtSubpId` int NOT NULL,
  `iPartId` int NOT NULL,
  `iProcessId` int NOT NULL,
  `vSubPrgName` varchar(200) NOT NULL,
  `iCycleTime` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `planned_downtime`
--

CREATE TABLE `planned_downtime` (
  `iDownTimeId` int NOT NULL,
  `dDowntimeDate` date NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `product`
--

CREATE TABLE `product` (
  `id` int NOT NULL,
  `product_name` varchar(100) NOT NULL,
  `category_id` int NOT NULL,
  `vProductImage` varchar(50) NOT NULL DEFAULT 'no_image_available.png',
  `del_status` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `production_schedule`
--

CREATE TABLE `production_schedule` (
  `iSheduleId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `dScheduleDate` date NOT NULL,
  `dScheduleEndDate` date NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL,
  `eSchedule_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 - Single  Day, 2 - Multiple Day',
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `qms_data`
--

CREATE TABLE `qms_data` (
  `QmsId` int NOT NULL,
  `vPart` varchar(255) NOT NULL,
  `vDimention` varchar(100) NOT NULL,
  `dActual` decimal(10,2) NOT NULL,
  `dNominal` decimal(10,2) NOT NULL,
  `dLoTol` decimal(10,2) NOT NULL,
  `dHiTol` decimal(10,2) NOT NULL,
  `dDeviation` decimal(10,2) NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `report_utilization_data`
--

CREATE TABLE `report_utilization_data` (
  `iRudId` int NOT NULL,
  `iUserId` int DEFAULT NULL,
  `iMachineId` int DEFAULT NULL,
  `iModeId` int DEFAULT NULL,
  `iDataId` int DEFAULT NULL,
  `iShiftId` int DEFAULT NULL,
  `iTotSeconds` int DEFAULT NULL,
  `extraSeconds` int DEFAULT NULL,
  `startDate` datetime DEFAULT NULL,
  `endDate` datetime DEFAULT NULL,
  `dCDate` date DEFAULT NULL,
  `dSDate` date DEFAULT NULL,
  `dCreatedDate` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `rmt_customer`
--

CREATE TABLE `rmt_customer` (
  `iUserId` int NOT NULL,
  `vCompanyName` varchar(250) NOT NULL,
  `vRemotePath` varchar(250) NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `rmt_mm_machines`
--

CREATE TABLE `rmt_mm_machines` (
  `iMachineId` int NOT NULL,
  `iUserId` int NOT NULL,
  `iRemoteID` int NOT NULL,
  `vMachineName` varchar(250) NOT NULL,
  `vMachineModel` varchar(250) NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `schedule_items`
--

CREATE TABLE `schedule_items` (
  `iItemId` int NOT NULL,
  `iSheduleId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `iShiftId` int NOT NULL,
  `iPartId` int NOT NULL,
  `iProcessId` int NOT NULL,
  `iQty` int NOT NULL,
  `iOperatorId` int NOT NULL,
  `vSupervisorId` varchar(50) NOT NULL,
  `dSstartDate` date NOT NULL,
  `tStartTime` time NOT NULL,
  `tEndTime` time NOT NULL,
  `vJobId` varchar(50) NOT NULL,
  `dTotalHours` float NOT NULL,
  `eUlUsed` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0=> not Taken, 1=> Not  Taken '
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `seven_configration`
--

CREATE TABLE `seven_configration` (
  `iConfigId` int NOT NULL,
  `eConfigType` enum('Smart Operator Panel','Host Computer','Technical Support') NOT NULL,
  `iNoOfMachine` int NOT NULL,
  `iNoOfVisit` int NOT NULL,
  `iNoOfTraining` int NOT NULL,
  `dPrice` decimal(12,0) NOT NULL,
  `iCreatedBy` int NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `iLastBy` int NOT NULL,
  `dUpdatedDate` datetime NOT NULL,
  `iDeleteBy` int NOT NULL,
  `dDeletedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d','') NOT NULL DEFAULT 'y'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `shift`
--

CREATE TABLE `shift` (
  `iShiftId` int NOT NULL,
  `iUserId` int NOT NULL,
  `vShiftName` varchar(100) NOT NULL,
  `eShiftType` enum('Day','Night') NOT NULL,
  `dShiftHours` double(5,3) NOT NULL,
  `tStartTime` time NOT NULL,
  `tEndTime` time NOT NULL,
  `vBreakId` varchar(20) NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `sub_licencekeys`
--

CREATE TABLE `sub_licencekeys` (
  `iLicencekeyId` int NOT NULL,
  `vLicenceKey` varchar(50) NOT NULL,
  `iOrdDetailId` int NOT NULL,
  `iPlanDuration` int NOT NULL,
  `eDurationType` enum('','Month','Year') NOT NULL,
  `iMachineId` int NOT NULL,
  `dActivationDate` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `sub_licencekeys_history`
--

CREATE TABLE `sub_licencekeys_history` (
  `iLicencekeyHisId` int NOT NULL,
  `iLicencekeyId` int NOT NULL,
  `vLicenceKey` varchar(50) NOT NULL,
  `iOrdDetailId` int NOT NULL,
  `iPlanDuration` int NOT NULL,
  `eDurationType` enum('','Month','Year') NOT NULL,
  `iMachineId` int NOT NULL,
  `dActivationDate` date NOT NULL,
  `dRenewDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `sub_orders`
--

CREATE TABLE `sub_orders` (
  `iOrderId` int NOT NULL,
  `iUserId` int NOT NULL,
  `iNoOfKey` int NOT NULL,
  `dOrderTotal` float NOT NULL,
  `dTaxAmount` float NOT NULL,
  `dFullTotal` float NOT NULL,
  `ePaymentMethod` enum('Online','Offline') NOT NULL,
  `ePaymentStatus` enum('Pending','Paid') NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `vPONumber` varchar(200) NOT NULL,
  `vPODate` date NOT NULL,
  `vPOAttach` varchar(100) NOT NULL,
  `vInvoiceFile` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `sub_order_details`
--

CREATE TABLE `sub_order_details` (
  `iOrdDetailId` int NOT NULL,
  `iOrderId` int NOT NULL,
  `vLicenceKey` varchar(50) NOT NULL,
  `iPlanId` int NOT NULL,
  `dPrice` float NOT NULL,
  `dTaxPrice` float NOT NULL,
  `dFullPrice` float NOT NULL,
  `eBindStatus` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 - unBinded 1 - Binded 2- Expired 3 - Renewed',
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `sub_payments`
--

CREATE TABLE `sub_payments` (
  `iPaymentId` int NOT NULL,
  `iUserId` int NOT NULL,
  `iPaymentMode` int NOT NULL,
  `vChequeNumber` varchar(100) NOT NULL,
  `dAmount` float NOT NULL,
  `vRemark` varchar(250) NOT NULL,
  `eStatus` enum('Pending','Received','d') NOT NULL,
  `dPaymentDate` date NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `dCreditedDate` date NOT NULL,
  `vConfirmationNo` varchar(200) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `sub_payments_transaction`
--

CREATE TABLE `sub_payments_transaction` (
  `iTransactionId` int NOT NULL,
  `iUserId` int NOT NULL,
  `iPaymentId` int NOT NULL COMMENT 'ONLY FOR CREDIT',
  `iOrderId` int NOT NULL COMMENT 'ONLY FOR DEBIT',
  `vTransActivity` varchar(200) NOT NULL,
  `vAmount` float NOT NULL,
  `eTransactionType` enum('Debit','Credit') NOT NULL,
  `fBalance` float NOT NULL,
  `vDetails` varchar(250) NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `sub_planduration`
--

CREATE TABLE `sub_planduration` (
  `iDurationId` int NOT NULL,
  `iDuration` varchar(10) NOT NULL,
  `eDurationType` enum('','Month','Year','Life Time','Considerable Trial','1 Year Considerable Trial','Smart Operator Panel','Siemens License') CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `sub_plans`
--

CREATE TABLE `sub_plans` (
  `iPlanId` int NOT NULL,
  `vPlanName` varchar(250) NOT NULL,
  `iDurationId` int NOT NULL COMMENT '0=Life Time',
  `iMinMachines` int NOT NULL,
  `iMaxMachines` int NOT NULL COMMENT '1000=UNLIMITED',
  `dPrice` float NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tbl_email_template`
--

CREATE TABLE `tbl_email_template` (
  `iTemplateID` int UNSIGNED NOT NULL,
  `vSubject` varchar(255) NOT NULL,
  `tEmailBody` text NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- --------------------------------------------------------

--
-- Table structure for table `tbl_users`
--

CREATE TABLE `tbl_users` (
  `iUserId` bigint UNSIGNED NOT NULL,
  `iRoleId` int NOT NULL,
  `iBranchId` int NOT NULL,
  `vUserName` varchar(100) NOT NULL,
  `vPassword` varchar(250) NOT NULL,
  `vFirstName` varchar(100) NOT NULL DEFAULT '' COMMENT 'COMPANY NAME',
  `vLastName` varchar(100) NOT NULL DEFAULT '',
  `iCompanyId` int NOT NULL,
  `vDesignation` varchar(100) NOT NULL,
  `vGrade` varchar(10) NOT NULL,
  `vMobileNumber` varchar(15) NOT NULL DEFAULT '',
  `dBirthDate` date NOT NULL,
  `dJoiningDate` date NOT NULL,
  `vBloodGroup` varchar(10) NOT NULL,
  `vProfilePicture` tinytext NOT NULL,
  `vEmail` tinytext NOT NULL,
  `vContactPerson` varchar(40) NOT NULL,
  `vScPhone` varchar(40) NOT NULL,
  `iCountryId` int NOT NULL,
  `iStateId` int NOT NULL,
  `vCity` varchar(45) NOT NULL,
  `vAddress1` varchar(250) NOT NULL,
  `vAddress2` varchar(250) NOT NULL,
  `vPinCode` varchar(15) NOT NULL,
  `eStatus` enum('y','n','d') NOT NULL,
  `eTheme` enum('dark_version','light_version') NOT NULL,
  `vTempCode` varchar(25) NOT NULL,
  `dOTPDateTime` datetime NOT NULL,
  `dLastActivity` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `dExpireDate` date NOT NULL,
  `dCreatedDate` datetime NOT NULL,
  `iFinixId` int NOT NULL,
  `vGstNo` varchar(20) NOT NULL,
  `vPanNumber` varchar(20) NOT NULL,
  `vWebSite` varchar(100) NOT NULL,
  `aMachineIds` text NOT NULL,
  `vTimeZone` varchar(255) NOT NULL,
  `tStartDayTime` time NOT NULL,
  `vLang` varchar(10) NOT NULL,
  `api_token` varchar(255) NOT NULL,
  `dReportCronDate` date NOT NULL COMMENT 'calcnew\\cron_report_data.php',
  `dLiveMonitorCronDate` date NOT NULL COMMENT 'calcnew\\cron_mm_data.php'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

-- --------------------------------------------------------

--
-- Table structure for table `tbl_users_macaddress`
--

CREATE TABLE `tbl_users_macaddress` (
  `iUsMacId` int NOT NULL,
  `iUserId` int NOT NULL,
  `vMacAddress` varchar(100) NOT NULL,
  `dLastActivity` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tbl_users_role`
--

CREATE TABLE `tbl_users_role` (
  `iRoleId` int NOT NULL,
  `vRoleName` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tbl_users_setting`
--

CREATE TABLE `tbl_users_setting` (
  `iSettingId` int NOT NULL,
  `iUserId` int NOT NULL,
  `tMinDtNtfcaion` enum('Yes','No') NOT NULL DEFAULT 'No',
  `tMinDwnTime` time NOT NULL,
  `tDwnNtfction` enum('Yes','No') NOT NULL DEFAULT 'No',
  `tDwnNtfctionTime` time NOT NULL,
  `eReportNtfcion` enum('Yes','No') NOT NULL DEFAULT 'No',
  `eReportPeriod` enum('Daily','Weekly','Monthly') NOT NULL,
  `tReportNtfcionTime` time NOT NULL,
  `vNtfcionReportIds` varchar(100) NOT NULL,
  `ePartTolerance` enum('Yes','No') NOT NULL DEFAULT 'No',
  `tPartTolerancePer` decimal(10,2) NOT NULL,
  `eCustomTime` enum('Yes','No') NOT NULL,
  `tCustomeTime` time NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ticketcenter`
--

CREATE TABLE `ticketcenter` (
  `iTcId` int NOT NULL,
  `iUserId` int NOT NULL,
  `vUname` varchar(255) NOT NULL,
  `vUemail` varchar(255) NOT NULL,
  `vUPhone` varchar(40) NOT NULL,
  `vSubject` varchar(255) NOT NULL,
  `tMessage` text NOT NULL,
  `iReplyTcId` int NOT NULL,
  `eStatus` enum('Pending','In Progress','Closed','Invalid','d') NOT NULL DEFAULT 'Pending',
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ticketcenter_attachment`
--

CREATE TABLE `ticketcenter_attachment` (
  `iTcAtchId` int NOT NULL,
  `iTcId` int NOT NULL,
  `vFile` varchar(255) NOT NULL,
  `vFileName` varchar(255) NOT NULL,
  `dCreatedDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `weekoff`
--

CREATE TABLE `weekoff` (
  `iWeekOffId` int NOT NULL,
  `iUserId` int NOT NULL,
  `vHolidayName` varchar(100) DEFAULT NULL,
  `eWeekOffType` enum('1','2') NOT NULL COMMENT '1=holiday, 2=weekoff',
  `dStartDate` date NOT NULL,
  `dEndDate` date NOT NULL,
  `iDay` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `zz_notifications`
--

CREATE TABLE `zz_notifications` (
  `iNotificationId` int NOT NULL,
  `iMachineId` int NOT NULL,
  `vMessage` varchar(250) NOT NULL,
  `isRead` tinyint(1) NOT NULL,
  `dCreatedDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `zz_timzone`
--

CREATE TABLE `zz_timzone` (
  `iTzId` int NOT NULL,
  `vTzName` varchar(200) NOT NULL,
  `vTzOffset` varchar(10) DEFAULT NULL,
  `vTzStandard` varchar(100) DEFAULT NULL,
  `eStatus` enum('Y','N') NOT NULL DEFAULT 'Y'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `zz_translate`
--

CREATE TABLE `zz_translate` (
  `tId` int NOT NULL,
  `vKeyword` varchar(200) NOT NULL,
  `vEnglish` varchar(200) NOT NULL,
  `vGerman` varchar(200) NOT NULL,
  `vFrench` varchar(200) NOT NULL,
  `vItalian` varchar(200) NOT NULL,
  `dLastDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `batch_tbl`
--
ALTER TABLE `batch_tbl`
  ADD PRIMARY KEY (`batchId`);

--
-- Indexes for table `branches`
--
ALTER TABLE `branches`
  ADD PRIMARY KEY (`iBranchId`),
  ADD KEY `id` (`iBranchId`),
  ADD KEY `iBranchGroup` (`iBranchGroup`);

--
-- Indexes for table `brand`
--
ALTER TABLE `brand`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `breaks`
--
ALTER TABLE `breaks`
  ADD PRIMARY KEY (`iBreakId`),
  ADD UNIQUE KEY `iId` (`iBreakId`);

--
-- Indexes for table `category`
--
ALTER TABLE `category`
  ADD PRIMARY KEY (`id`),
  ADD KEY `iMCategoryId` (`iMCategoryId`),
  ADD KEY `id` (`id`),
  ADD KEY `id_2` (`id`),
  ADD KEY `iMCategoryId_2` (`iMCategoryId`);

--
-- Indexes for table `ci_sessions`
--
ALTER TABLE `ci_sessions`
  ADD KEY `ci_sessions_timestamp` (`timestamp`);

--
-- Indexes for table `clientele`
--
ALTER TABLE `clientele`
  ADD PRIMARY KEY (`iClienteleId`);

--
-- Indexes for table `controller`
--
ALTER TABLE `controller`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `downtime_items`
--
ALTER TABLE `downtime_items`
  ADD PRIMARY KEY (`iItemId`);

--
-- Indexes for table `downtime_reason`
--
ALTER TABLE `downtime_reason`
  ADD PRIMARY KEY (`iDtrId`);

--
-- Indexes for table `dropdown_key`
--
ALTER TABLE `dropdown_key`
  ADD PRIMARY KEY (`iKeyId`);

--
-- Indexes for table `dropdown_value`
--
ALTER TABLE `dropdown_value`
  ADD PRIMARY KEY (`iValueId`);

--
-- Indexes for table `helpenter_pages`
--
ALTER TABLE `helpenter_pages`
  ADD PRIMARY KEY (`iHcId`);

--
-- Indexes for table `helpenter_pages_bkup`
--
ALTER TABLE `helpenter_pages_bkup`
  ADD PRIMARY KEY (`iHcId`);

--
-- Indexes for table `idle_reason`
--
ALTER TABLE `idle_reason`
  ADD PRIMARY KEY (`iIdlerId`);

--
-- Indexes for table `machines_data`
--
ALTER TABLE `machines_data`
  ADD PRIMARY KEY (`iDataId`),
  ADD KEY `iMachineId` (`iMachineId`),
  ADD KEY `iModeId` (`iModeId`),
  ADD KEY `isUpdateFinix` (`isUpdateFinix`);

--
-- Indexes for table `machines_data_bkp`
--
ALTER TABLE `machines_data_bkp`
  ADD PRIMARY KEY (`iDataId`);

--
-- Indexes for table `machines_data_history`
--
ALTER TABLE `machines_data_history`
  ADD PRIMARY KEY (`iDataId`);

--
-- Indexes for table `machines_data_temp`
--
ALTER TABLE `machines_data_temp`
  ADD PRIMARY KEY (`iDataId`),
  ADD KEY `iMachineId` (`iMachineId`),
  ADD KEY `iModeId` (`iModeId`),
  ADD KEY `isUpdateFinix` (`isUpdateFinix`);

--
-- Indexes for table `machines_modes`
--
ALTER TABLE `machines_modes`
  ADD PRIMARY KEY (`iModeId`);

--
-- Indexes for table `machine_accuracy`
--
ALTER TABLE `machine_accuracy`
  ADD UNIQUE KEY `iId` (`iAccuracyId`);

--
-- Indexes for table `machine_alarms`
--
ALTER TABLE `machine_alarms`
  ADD PRIMARY KEY (`iAlarmId`);

--
-- Indexes for table `machine_alarms_history`
--
ALTER TABLE `machine_alarms_history`
  ADD PRIMARY KEY (`iAlarmId`);

--
-- Indexes for table `machine_data_count`
--
ALTER TABLE `machine_data_count`
  ADD PRIMARY KEY (`iMdcId`);

--
-- Indexes for table `machine_data_count_history`
--
ALTER TABLE `machine_data_count_history`
  ADD PRIMARY KEY (`iMdcId`);

--
-- Indexes for table `machine_data_count_hourly`
--
ALTER TABLE `machine_data_count_hourly`
  ADD PRIMARY KEY (`iMdcId`);

--
-- Indexes for table `machine_data_count_hourly_history`
--
ALTER TABLE `machine_data_count_hourly_history`
  ADD PRIMARY KEY (`iMdcId`);

--
-- Indexes for table `machine_lineno`
--
ALTER TABLE `machine_lineno`
  ADD UNIQUE KEY `iId` (`iLineId`);

--
-- Indexes for table `machine_online_log`
--
ALTER TABLE `machine_online_log`
  ADD PRIMARY KEY (`iOnLogid`);

--
-- Indexes for table `machine_programs`
--
ALTER TABLE `machine_programs`
  ADD PRIMARY KEY (`iProgramId`);

--
-- Indexes for table `maintenance_tasks`
--
ALTER TABLE `maintenance_tasks`
  ADD PRIMARY KEY (`iMTaskId`);

--
-- Indexes for table `mm_api_timelog`
--
ALTER TABLE `mm_api_timelog`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `mm_cycletime`
--
ALTER TABLE `mm_cycletime`
  ADD PRIMARY KEY (`iDataId`);

--
-- Indexes for table `mm_cycletime_history`
--
ALTER TABLE `mm_cycletime_history`
  ADD PRIMARY KEY (`iDataId`);

--
-- Indexes for table `mm_daily_production`
--
ALTER TABLE `mm_daily_production`
  ADD PRIMARY KEY (`iDlyPrdId`);

--
-- Indexes for table `mm_daily_production_history`
--
ALTER TABLE `mm_daily_production_history`
  ADD PRIMARY KEY (`iDphId`);

--
-- Indexes for table `mm_health`
--
ALTER TABLE `mm_health`
  ADD PRIMARY KEY (`iHealthId`);

--
-- Indexes for table `mm_machines`
--
ALTER TABLE `mm_machines`
  ADD PRIMARY KEY (`iMachineId`);

--
-- Indexes for table `mm_machines_bkp`
--
ALTER TABLE `mm_machines_bkp`
  ADD PRIMARY KEY (`iMachineId`);

--
-- Indexes for table `mm_machines_ori`
--
ALTER TABLE `mm_machines_ori`
  ADD PRIMARY KEY (`iMachineId`);

--
-- Indexes for table `mm_modules`
--
ALTER TABLE `mm_modules`
  ADD PRIMARY KEY (`iModuleId`);

--
-- Indexes for table `mm_modules_old`
--
ALTER TABLE `mm_modules_old`
  ADD PRIMARY KEY (`iModuleId`);

--
-- Indexes for table `mm_operator_tasks`
--
ALTER TABLE `mm_operator_tasks`
  ADD PRIMARY KEY (`iOpTaskId`);

--
-- Indexes for table `mm_programs`
--
ALTER TABLE `mm_programs`
  ADD PRIMARY KEY (`iPrgId`);

--
-- Indexes for table `mm_programs_history`
--
ALTER TABLE `mm_programs_history`
  ADD PRIMARY KEY (`iPrgId`);

--
-- Indexes for table `mm_tool_life`
--
ALTER TABLE `mm_tool_life`
  ADD PRIMARY KEY (`iTlifeId`);

--
-- Indexes for table `mm_tool_life_history`
--
ALTER TABLE `mm_tool_life_history`
  ADD PRIMARY KEY (`iTlifeId`);

--
-- Indexes for table `mm_variables`
--
ALTER TABLE `mm_variables`
  ADD PRIMARY KEY (`iVarId`);

--
-- Indexes for table `mm_variables_01-11-2022`
--
ALTER TABLE `mm_variables_01-11-2022`
  ADD PRIMARY KEY (`iVarId`);

--
-- Indexes for table `mm_variables_21-09-2022`
--
ALTER TABLE `mm_variables_21-09-2022`
  ADD PRIMARY KEY (`iVarId`);

--
-- Indexes for table `mm_variables_final`
--
ALTER TABLE `mm_variables_final`
  ADD PRIMARY KEY (`iVarId`);

--
-- Indexes for table `module_permission`
--
ALTER TABLE `module_permission`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `mst_countries`
--
ALTER TABLE `mst_countries`
  ADD PRIMARY KEY (`iCountryId`);

--
-- Indexes for table `mst_downtime_reason`
--
ALTER TABLE `mst_downtime_reason`
  ADD PRIMARY KEY (`iDownReason`);

--
-- Indexes for table `mst_downtime_reason_copy`
--
ALTER TABLE `mst_downtime_reason_copy`
  ADD PRIMARY KEY (`iDownReason`);

--
-- Indexes for table `mst_idle_reason`
--
ALTER TABLE `mst_idle_reason`
  ADD PRIMARY KEY (`iIdleReason`);

--
-- Indexes for table `mst_imagethumb`
--
ALTER TABLE `mst_imagethumb`
  ADD PRIMARY KEY (`iThumbId`);

--
-- Indexes for table `mst_site_settings`
--
ALTER TABLE `mst_site_settings`
  ADD PRIMARY KEY (`iFieldId`);

--
-- Indexes for table `mst_states`
--
ALTER TABLE `mst_states`
  ADD PRIMARY KEY (`iStateId`);

--
-- Indexes for table `oper_mainten_tasks`
--
ALTER TABLE `oper_mainten_tasks`
  ADD PRIMARY KEY (`iOpMtskId`);

--
-- Indexes for table `parts`
--
ALTER TABLE `parts`
  ADD PRIMARY KEY (`iPartId`);

--
-- Indexes for table `parts_process`
--
ALTER TABLE `parts_process`
  ADD PRIMARY KEY (`iProcessId`);

--
-- Indexes for table `parts_sub_process`
--
ALTER TABLE `parts_sub_process`
  ADD PRIMARY KEY (`iPrtSubpId`);

--
-- Indexes for table `planned_downtime`
--
ALTER TABLE `planned_downtime`
  ADD PRIMARY KEY (`iDownTimeId`);

--
-- Indexes for table `product`
--
ALTER TABLE `product`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id` (`id`),
  ADD KEY `id_2` (`id`);

--
-- Indexes for table `production_schedule`
--
ALTER TABLE `production_schedule`
  ADD PRIMARY KEY (`iSheduleId`);

--
-- Indexes for table `qms_data`
--
ALTER TABLE `qms_data`
  ADD PRIMARY KEY (`QmsId`);

--
-- Indexes for table `report_utilization_data`
--
ALTER TABLE `report_utilization_data`
  ADD PRIMARY KEY (`iRudId`);

--
-- Indexes for table `rmt_customer`
--
ALTER TABLE `rmt_customer`
  ADD PRIMARY KEY (`iUserId`);

--
-- Indexes for table `rmt_mm_machines`
--
ALTER TABLE `rmt_mm_machines`
  ADD PRIMARY KEY (`iMachineId`);

--
-- Indexes for table `schedule_items`
--
ALTER TABLE `schedule_items`
  ADD PRIMARY KEY (`iItemId`);

--
-- Indexes for table `seven_configration`
--
ALTER TABLE `seven_configration`
  ADD PRIMARY KEY (`iConfigId`);

--
-- Indexes for table `shift`
--
ALTER TABLE `shift`
  ADD PRIMARY KEY (`iShiftId`),
  ADD UNIQUE KEY `iId` (`iShiftId`);

--
-- Indexes for table `sub_licencekeys`
--
ALTER TABLE `sub_licencekeys`
  ADD PRIMARY KEY (`iLicencekeyId`);

--
-- Indexes for table `sub_licencekeys_history`
--
ALTER TABLE `sub_licencekeys_history`
  ADD PRIMARY KEY (`iLicencekeyHisId`);

--
-- Indexes for table `sub_orders`
--
ALTER TABLE `sub_orders`
  ADD PRIMARY KEY (`iOrderId`);

--
-- Indexes for table `sub_order_details`
--
ALTER TABLE `sub_order_details`
  ADD PRIMARY KEY (`iOrdDetailId`);

--
-- Indexes for table `sub_payments`
--
ALTER TABLE `sub_payments`
  ADD PRIMARY KEY (`iPaymentId`);

--
-- Indexes for table `sub_payments_transaction`
--
ALTER TABLE `sub_payments_transaction`
  ADD PRIMARY KEY (`iTransactionId`);

--
-- Indexes for table `sub_planduration`
--
ALTER TABLE `sub_planduration`
  ADD PRIMARY KEY (`iDurationId`);

--
-- Indexes for table `sub_plans`
--
ALTER TABLE `sub_plans`
  ADD PRIMARY KEY (`iPlanId`);

--
-- Indexes for table `tbl_email_template`
--
ALTER TABLE `tbl_email_template`
  ADD PRIMARY KEY (`iTemplateID`);

--
-- Indexes for table `tbl_users`
--
ALTER TABLE `tbl_users`
  ADD PRIMARY KEY (`iUserId`);

--
-- Indexes for table `tbl_users_macaddress`
--
ALTER TABLE `tbl_users_macaddress`
  ADD PRIMARY KEY (`iUsMacId`);

--
-- Indexes for table `tbl_users_role`
--
ALTER TABLE `tbl_users_role`
  ADD PRIMARY KEY (`iRoleId`);

--
-- Indexes for table `tbl_users_setting`
--
ALTER TABLE `tbl_users_setting`
  ADD PRIMARY KEY (`iSettingId`);

--
-- Indexes for table `ticketcenter`
--
ALTER TABLE `ticketcenter`
  ADD PRIMARY KEY (`iTcId`);

--
-- Indexes for table `ticketcenter_attachment`
--
ALTER TABLE `ticketcenter_attachment`
  ADD PRIMARY KEY (`iTcAtchId`);

--
-- Indexes for table `weekoff`
--
ALTER TABLE `weekoff`
  ADD PRIMARY KEY (`iWeekOffId`);

--
-- Indexes for table `zz_notifications`
--
ALTER TABLE `zz_notifications`
  ADD PRIMARY KEY (`iNotificationId`);

--
-- Indexes for table `zz_timzone`
--
ALTER TABLE `zz_timzone`
  ADD PRIMARY KEY (`iTzId`);

--
-- Indexes for table `zz_translate`
--
ALTER TABLE `zz_translate`
  ADD PRIMARY KEY (`tId`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `batch_tbl`
--
ALTER TABLE `batch_tbl`
  MODIFY `batchId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `branches`
--
ALTER TABLE `branches`
  MODIFY `iBranchId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `brand`
--
ALTER TABLE `brand`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `breaks`
--
ALTER TABLE `breaks`
  MODIFY `iBreakId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `category`
--
ALTER TABLE `category`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `clientele`
--
ALTER TABLE `clientele`
  MODIFY `iClienteleId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `controller`
--
ALTER TABLE `controller`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `downtime_items`
--
ALTER TABLE `downtime_items`
  MODIFY `iItemId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `downtime_reason`
--
ALTER TABLE `downtime_reason`
  MODIFY `iDtrId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `dropdown_key`
--
ALTER TABLE `dropdown_key`
  MODIFY `iKeyId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `dropdown_value`
--
ALTER TABLE `dropdown_value`
  MODIFY `iValueId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `helpenter_pages`
--
ALTER TABLE `helpenter_pages`
  MODIFY `iHcId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `helpenter_pages_bkup`
--
ALTER TABLE `helpenter_pages_bkup`
  MODIFY `iHcId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `idle_reason`
--
ALTER TABLE `idle_reason`
  MODIFY `iIdlerId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machines_data`
--
ALTER TABLE `machines_data`
  MODIFY `iDataId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machines_data_bkp`
--
ALTER TABLE `machines_data_bkp`
  MODIFY `iDataId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machines_data_history`
--
ALTER TABLE `machines_data_history`
  MODIFY `iDataId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machines_data_temp`
--
ALTER TABLE `machines_data_temp`
  MODIFY `iDataId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machines_modes`
--
ALTER TABLE `machines_modes`
  MODIFY `iModeId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_accuracy`
--
ALTER TABLE `machine_accuracy`
  MODIFY `iAccuracyId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_alarms`
--
ALTER TABLE `machine_alarms`
  MODIFY `iAlarmId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_alarms_history`
--
ALTER TABLE `machine_alarms_history`
  MODIFY `iAlarmId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_data_count`
--
ALTER TABLE `machine_data_count`
  MODIFY `iMdcId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_data_count_history`
--
ALTER TABLE `machine_data_count_history`
  MODIFY `iMdcId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_data_count_hourly`
--
ALTER TABLE `machine_data_count_hourly`
  MODIFY `iMdcId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_data_count_hourly_history`
--
ALTER TABLE `machine_data_count_hourly_history`
  MODIFY `iMdcId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_lineno`
--
ALTER TABLE `machine_lineno`
  MODIFY `iLineId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_online_log`
--
ALTER TABLE `machine_online_log`
  MODIFY `iOnLogid` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `machine_programs`
--
ALTER TABLE `machine_programs`
  MODIFY `iProgramId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `maintenance_tasks`
--
ALTER TABLE `maintenance_tasks`
  MODIFY `iMTaskId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_api_timelog`
--
ALTER TABLE `mm_api_timelog`
  MODIFY `Id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_cycletime`
--
ALTER TABLE `mm_cycletime`
  MODIFY `iDataId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_cycletime_history`
--
ALTER TABLE `mm_cycletime_history`
  MODIFY `iDataId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_daily_production`
--
ALTER TABLE `mm_daily_production`
  MODIFY `iDlyPrdId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_daily_production_history`
--
ALTER TABLE `mm_daily_production_history`
  MODIFY `iDphId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_health`
--
ALTER TABLE `mm_health`
  MODIFY `iHealthId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_machines`
--
ALTER TABLE `mm_machines`
  MODIFY `iMachineId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_machines_bkp`
--
ALTER TABLE `mm_machines_bkp`
  MODIFY `iMachineId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_machines_ori`
--
ALTER TABLE `mm_machines_ori`
  MODIFY `iMachineId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_modules`
--
ALTER TABLE `mm_modules`
  MODIFY `iModuleId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_modules_old`
--
ALTER TABLE `mm_modules_old`
  MODIFY `iModuleId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_operator_tasks`
--
ALTER TABLE `mm_operator_tasks`
  MODIFY `iOpTaskId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_programs`
--
ALTER TABLE `mm_programs`
  MODIFY `iPrgId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_programs_history`
--
ALTER TABLE `mm_programs_history`
  MODIFY `iPrgId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_tool_life`
--
ALTER TABLE `mm_tool_life`
  MODIFY `iTlifeId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_tool_life_history`
--
ALTER TABLE `mm_tool_life_history`
  MODIFY `iTlifeId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_variables`
--
ALTER TABLE `mm_variables`
  MODIFY `iVarId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_variables_01-11-2022`
--
ALTER TABLE `mm_variables_01-11-2022`
  MODIFY `iVarId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_variables_21-09-2022`
--
ALTER TABLE `mm_variables_21-09-2022`
  MODIFY `iVarId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mm_variables_final`
--
ALTER TABLE `mm_variables_final`
  MODIFY `iVarId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `module_permission`
--
ALTER TABLE `module_permission`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mst_countries`
--
ALTER TABLE `mst_countries`
  MODIFY `iCountryId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mst_downtime_reason`
--
ALTER TABLE `mst_downtime_reason`
  MODIFY `iDownReason` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mst_downtime_reason_copy`
--
ALTER TABLE `mst_downtime_reason_copy`
  MODIFY `iDownReason` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mst_idle_reason`
--
ALTER TABLE `mst_idle_reason`
  MODIFY `iIdleReason` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mst_imagethumb`
--
ALTER TABLE `mst_imagethumb`
  MODIFY `iThumbId` int UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mst_site_settings`
--
ALTER TABLE `mst_site_settings`
  MODIFY `iFieldId` int UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mst_states`
--
ALTER TABLE `mst_states`
  MODIFY `iStateId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `oper_mainten_tasks`
--
ALTER TABLE `oper_mainten_tasks`
  MODIFY `iOpMtskId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `parts`
--
ALTER TABLE `parts`
  MODIFY `iPartId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `parts_process`
--
ALTER TABLE `parts_process`
  MODIFY `iProcessId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `parts_sub_process`
--
ALTER TABLE `parts_sub_process`
  MODIFY `iPrtSubpId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `planned_downtime`
--
ALTER TABLE `planned_downtime`
  MODIFY `iDownTimeId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `product`
--
ALTER TABLE `product`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `production_schedule`
--
ALTER TABLE `production_schedule`
  MODIFY `iSheduleId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `qms_data`
--
ALTER TABLE `qms_data`
  MODIFY `QmsId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `report_utilization_data`
--
ALTER TABLE `report_utilization_data`
  MODIFY `iRudId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `rmt_customer`
--
ALTER TABLE `rmt_customer`
  MODIFY `iUserId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `rmt_mm_machines`
--
ALTER TABLE `rmt_mm_machines`
  MODIFY `iMachineId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `schedule_items`
--
ALTER TABLE `schedule_items`
  MODIFY `iItemId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `seven_configration`
--
ALTER TABLE `seven_configration`
  MODIFY `iConfigId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `shift`
--
ALTER TABLE `shift`
  MODIFY `iShiftId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sub_licencekeys`
--
ALTER TABLE `sub_licencekeys`
  MODIFY `iLicencekeyId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sub_licencekeys_history`
--
ALTER TABLE `sub_licencekeys_history`
  MODIFY `iLicencekeyHisId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sub_orders`
--
ALTER TABLE `sub_orders`
  MODIFY `iOrderId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sub_order_details`
--
ALTER TABLE `sub_order_details`
  MODIFY `iOrdDetailId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sub_payments`
--
ALTER TABLE `sub_payments`
  MODIFY `iPaymentId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sub_payments_transaction`
--
ALTER TABLE `sub_payments_transaction`
  MODIFY `iTransactionId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sub_planduration`
--
ALTER TABLE `sub_planduration`
  MODIFY `iDurationId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sub_plans`
--
ALTER TABLE `sub_plans`
  MODIFY `iPlanId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tbl_email_template`
--
ALTER TABLE `tbl_email_template`
  MODIFY `iTemplateID` int UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tbl_users`
--
ALTER TABLE `tbl_users`
  MODIFY `iUserId` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tbl_users_macaddress`
--
ALTER TABLE `tbl_users_macaddress`
  MODIFY `iUsMacId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tbl_users_role`
--
ALTER TABLE `tbl_users_role`
  MODIFY `iRoleId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tbl_users_setting`
--
ALTER TABLE `tbl_users_setting`
  MODIFY `iSettingId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ticketcenter`
--
ALTER TABLE `ticketcenter`
  MODIFY `iTcId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ticketcenter_attachment`
--
ALTER TABLE `ticketcenter_attachment`
  MODIFY `iTcAtchId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `weekoff`
--
ALTER TABLE `weekoff`
  MODIFY `iWeekOffId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `zz_notifications`
--
ALTER TABLE `zz_notifications`
  MODIFY `iNotificationId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `zz_timzone`
--
ALTER TABLE `zz_timzone`
  MODIFY `iTzId` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `zz_translate`
--
ALTER TABLE `zz_translate`
  MODIFY `tId` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
