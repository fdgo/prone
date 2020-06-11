/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.163.133
 Source Server Type    : MySQL
 Source Server Version : 80016
 Source Host           : 192.168.163.133:3306
 Source Schema         : go

 Target Server Type    : MySQL
 Target Server Version : 80016
 File Encoding         : 65001

 Date: 24/05/2020 16:19:50
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_menu
-- ----------------------------
DROP TABLE IF EXISTS `t_menu`;
CREATE TABLE `t_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '名字',
  `path` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '访问路径',
  `method` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '资源请求方式',
  `created_at` int(11) UNSIGNED DEFAULT NULL COMMENT '创建时间',
  `modified_at` int(11) UNSIGNED DEFAULT NULL COMMENT '更新时间',
  `deleted_at` int(11) UNSIGNED DEFAULT 0 COMMENT '删除时间戳',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of t_menu
-- ----------------------------
INSERT INTO `t_menu` VALUES (1, '查询所有菜单', '/api/v1/menus', 'GET', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (2, '查询单个菜单', '/api/v1/menus/:id', 'GET', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (3, '创建单个菜单', '/api/v1/menus', 'POST', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (4, '更新单个菜单', '/api/v1/menus/:id', 'PUT', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (5, '删除单个菜单', '/api/v1/menus/:id', 'DELETE', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (6, '查询所有用户', '/api/v1/users', 'GET', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (7, '查询单个用户', '/api/v1/users/:id', 'GET', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (8, '创建单个用户', '/api/v1/users', 'POST', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (9, '更新单个用户', '/api/v1/users/:id', 'PUT', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (10, '删除单个用户', '/api/v1/users/:id', 'DELETE', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (11, '查询所有角色', '/api/v1/roles', 'GET', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (12, '查询单个角色', '/api/v1/roles/:id', 'GET', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (13, '创建单个角色', '/api/v1/roles', 'POST', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (14, '更新单个角色', '/api/v1/roles/:id', 'PUT', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (15, '删除单个角色', '/api/v1/roles/:id', 'DELETE', NULL, NULL, 0);
INSERT INTO `t_menu` VALUES (16, '登录', '/auth', 'POST', NULL, NULL, 0);

-- ----------------------------
-- Table structure for t_role
-- ----------------------------
DROP TABLE IF EXISTS `t_role`;
CREATE TABLE `t_role`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '名字',
  `created_at` int(11) UNSIGNED DEFAULT NULL COMMENT '创建时间',
  `modified_at` int(11) UNSIGNED DEFAULT NULL COMMENT '更新时间',
  `deleted_at` int(11) UNSIGNED DEFAULT 0 COMMENT '删除时间戳',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of t_role
-- ----------------------------
INSERT INTO `t_role` VALUES (1, '开发部', NULL, NULL, 0);
INSERT INTO `t_role` VALUES (2, '运维部', NULL, NULL, 0);
INSERT INTO `t_role` VALUES (3, '测试部', NULL, NULL, 0);

-- ----------------------------
-- Table structure for t_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `t_role_menu`;
CREATE TABLE `t_role_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `role_id` int(11) UNSIGNED DEFAULT NULL COMMENT '角色ID',
  `menu_id` int(11) UNSIGNED DEFAULT NULL COMMENT '菜单ID',
  `deleted_at` int(11) UNSIGNED DEFAULT 0 COMMENT '删除时间戳',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户_角色ID_管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of t_role_menu
-- ----------------------------
INSERT INTO `t_role_menu` VALUES (1, 2, 1, 0);
INSERT INTO `t_role_menu` VALUES (2, 2, 2, 0);
INSERT INTO `t_role_menu` VALUES (3, 2, 3, 0);
INSERT INTO `t_role_menu` VALUES (4, 2, 4, 0);
INSERT INTO `t_role_menu` VALUES (5, 2, 5, 0);
INSERT INTO `t_role_menu` VALUES (6, 2, 6, 0);
INSERT INTO `t_role_menu` VALUES (7, 2, 7, 0);
INSERT INTO `t_role_menu` VALUES (8, 2, 8, 0);
INSERT INTO `t_role_menu` VALUES (9, 2, 9, 0);
INSERT INTO `t_role_menu` VALUES (10, 2, 10, 0);
INSERT INTO `t_role_menu` VALUES (11, 2, 11, 0);
INSERT INTO `t_role_menu` VALUES (12, 2, 12, 0);
INSERT INTO `t_role_menu` VALUES (13, 2, 13, 0);
INSERT INTO `t_role_menu` VALUES (14, 2, 14, 0);
INSERT INTO `t_role_menu` VALUES (15, 2, 15, 0);

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '账号',
  `password` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '密码',
  `created_at` int(11) UNSIGNED DEFAULT NULL COMMENT '创建时间',
  `modified_at` int(11) UNSIGNED DEFAULT NULL COMMENT '更新时间',
  `deleted_at` int(11) UNSIGNED DEFAULT 0 COMMENT '删除时间戳',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of t_user
-- ----------------------------
INSERT INTO `t_user` VALUES (1, 'admin', 'e10adc3949ba59abbe56e057f20f883e', NULL, NULL, 0);
INSERT INTO `t_user` VALUES (2, 'wmh', '14e1b600b1fd579f47433b88e8d85291', 1550642309, 1550642309, 0);

-- ----------------------------
-- Table structure for t_user_role
-- ----------------------------
DROP TABLE IF EXISTS `t_user_role`;
CREATE TABLE `t_user_role`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(11) UNSIGNED DEFAULT NULL COMMENT '用户ID',
  `role_id` int(11) UNSIGNED DEFAULT NULL COMMENT '角色ID',
  `deleted_at` int(11) UNSIGNED DEFAULT 0 COMMENT '删除时间戳',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户_角色ID_管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of t_user_role
-- ----------------------------
INSERT INTO `t_user_role` VALUES (1, 2, 2, 0);

SET FOREIGN_KEY_CHECKS = 1;
