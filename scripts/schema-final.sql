-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema auth
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema auth
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `auth` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ;
USE `auth` ;

-- -----------------------------------------------------
-- Table `auth`.`User`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `auth`.`User` ;

CREATE TABLE IF NOT EXISTS `auth`.`User` (
    `Id` CHAR(36) NOT NULL,
    `Email` VARCHAR(100) NOT NULL,
    `Password` VARCHAR(600) NOT NULL,
    `Status` ENUM('active', 'disabled', 'pending') NOT NULL,
    PRIMARY KEY (`Id`),
    UNIQUE INDEX `Email_UNIQUE` (`Email` ASC) VISIBLE)
    ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `auth`.`Verification`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `auth`.`Verification` ;

CREATE TABLE IF NOT EXISTS `auth`.`Verification` (
    `Code` VARCHAR(100) NOT NULL,
    `Type` ENUM('signup', 'recover') NOT NULL,
    `User_Id` CHAR(36) NOT NULL,
    UNIQUE INDEX `Code_UNIQUE` (`Code` ASC) VISIBLE,
    PRIMARY KEY (`User_Id`, `Type`),
    CONSTRAINT `fk_Verification_User`
    FOREIGN KEY (`User_Id`)
    REFERENCES `auth`.`User` (`Id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
    ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
