/*
 Navicat Premium Data Transfer

 Source Server         : local database
 Source Server Type    : MySQL
 Source Server Version : 80023
 Source Host           : localhost:3306
 Source Schema         : classroom_booking

 Target Server Type    : MySQL
 Target Server Version : 80023
 File Encoding         : 65001

 Date: 25/08/2023 14:58:18
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for bookings
-- ----------------------------
DROP TABLE IF EXISTS `bookings`;
CREATE TABLE `bookings`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `room_id` int NOT NULL,
  `user_id` int NOT NULL,
  `date` date NOT NULL,
  `start_time` time NOT NULL,
  `end_time` time NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `room_id`(`room_id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  CONSTRAINT `bookings_ibfk_1` FOREIGN KEY (`room_id`) REFERENCES `classrooms` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `bookings_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bookings
-- ----------------------------
INSERT INTO `bookings` VALUES (3, 4, 5, '2023-08-22', '09:00:00', '11:00:00');
INSERT INTO `bookings` VALUES (7, 4, 5, '2023-08-22', '12:00:00', '14:00:00');
INSERT INTO `bookings` VALUES (8, 4, 5, '2023-08-22', '15:00:00', '17:00:00');
INSERT INTO `bookings` VALUES (11, 4, 5, '2023-08-23', '12:00:00', '14:00:00');
INSERT INTO `bookings` VALUES (12, 5, 5, '2023-08-23', '12:00:00', '14:00:00');
INSERT INTO `bookings` VALUES (13, 8, 5, '2023-08-23', '12:00:00', '14:00:00');

-- ----------------------------
-- Table structure for classrooms
-- ----------------------------
DROP TABLE IF EXISTS `classrooms`;
CREATE TABLE `classrooms`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `floor` tinyint NOT NULL,
  `power` tinyint NOT NULL COMMENT '0表示不带电源，1表示带电源',
  `capacity` int NOT NULL,
  `photo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `roomName` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 19 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of classrooms
-- ----------------------------
INSERT INTO `classrooms` VALUES (0, 'ucl', 1, 0, 4, './images/3737366666393835383430643332346265663264656234383032396131366135.jpg', '101');
INSERT INTO `classrooms` VALUES (4, 'ucl', 1, 1, 2, './images/3736643236346637326363626665353466656561306431393537336362313961.jpg', '102');
INSERT INTO `classrooms` VALUES (5, 'ucl', 2, 1, 2, './images/3461663662643136623434326232373534313966356535653736313538343861.jpg', '201');
INSERT INTO `classrooms` VALUES (6, 'ucl', 2, 0, 2, './images/3935303663616235356362346666303837366233363266323935353338643136.jpg', '202');
INSERT INTO `classrooms` VALUES (7, 'ucl', 3, 0, 10, './images/3037363633373264303234613264383230643432393634376632646333383864.jpg', '301');
INSERT INTO `classrooms` VALUES (8, 'ucl', 3, 0, 7, './images/6263356361373462393736303233663766313331616530636133323633636138.jpg', '302');
INSERT INTO `classrooms` VALUES (9, 'ucl', 5, 1, 10, './images/6236653961636366373435326263396339663332326536366534613563333463.jpg', '501');
INSERT INTO `classrooms` VALUES (13, 'ucl', 5, 1, 6, './images/3530633530343064313164643838646139616139323661346565613731366166.jpg', '502');
INSERT INTO `classrooms` VALUES (15, 'ucl', 6, 0, 10, './images/6230373538383336336532303563663533383937373039666463383964633233.jpg', '601');
INSERT INTO `classrooms` VALUES (18, 'ucl', 6, 1, 6, './images/6230373538383336336532303563663533383937373039666463383964633233.jpg', '601');

-- ----------------------------
-- Table structure for user_preferences
-- ----------------------------
DROP TABLE IF EXISTS `user_preferences`;
CREATE TABLE `user_preferences`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `room_id` int NULL DEFAULT NULL,
  `user_id` int NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `room_id`(`room_id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  CONSTRAINT `user_preferences_ibfk_1` FOREIGN KEY (`room_id`) REFERENCES `classrooms` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `user_preferences_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_preferences
-- ----------------------------
INSERT INTO `user_preferences` VALUES (1, 4, 5);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `user_type` tinyint NOT NULL DEFAULT 0 COMMENT '0代表普通用户，1代表管理员用户',
  `third_session` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `id_thirdsession`(`id`, `third_session`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (3, '123456789', '123456789@qq.com', 0, '72f03bb1d0914bf39ad5b3cb72ef03f1');
INSERT INTO `users` VALUES (4, 'root', 'admin@qq.com', 1, '5c643fcb1a2d4dcc8878fcfb82db73b7');
INSERT INTO `users` VALUES (5, '123456', '6666666@qq.com', 0, 'e7a0b2c29fe2421f8b8c113b34fbdff6');

SET FOREIGN_KEY_CHECKS = 1;
