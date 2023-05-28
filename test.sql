CREATE TABLE USERS(
   Id_USERS INT AUTO_INCREMENT,
   email VARCHAR(100)  NOT NULL,
   password CHAR(255)  NOT NULL,
   firstName VARCHAR(50)  NOT NULL,
   lastName VARCHAR(50)  NOT NULL,
   profilePicture VARCHAR(100)  DEFAULT 'default.jpg',
   isCreatedAt DATETIME NOT NULL DEFAULT NOW(),
   lastSeen DATETIME NOT NULL DEFAULT NOW(),
   isBlocked VARCHAR(100)  NOT NULL DEFAULT 'not blocked',
   PRIMARY KEY(Id_USERS),
   UNIQUE(email)
);

CREATE TABLE CLIENTS(
   Id_CLIENTS INT AUTO_INCREMENT,
   fidelityPoints INT DEFAULT 0,
   streetName VARCHAR(100) ,
   country VARCHAR(50) ,
   city VARCHAR(100) ,
   streetNumber VARCHAR(10) ,
   phoneNumber VARCHAR(25) ,
   Id_USERS INT NOT NULL,
   PRIMARY KEY(Id_CLIENTS),
   UNIQUE(Id_USERS),
   FOREIGN KEY(Id_USERS) REFERENCES USERS(Id_USERS)
);

CREATE TABLE MANAGERS(
   Id_MANAGERS INT AUTO_INCREMENT,
   isItemManager BOOLEAN,
   isClientManager BOOLEAN,
   isContractorManager BOOLEAN,
   isSuperAdmin BOOLEAN,
   Id_USERS INT NOT NULL,
   PRIMARY KEY(Id_MANAGERS),
   UNIQUE(Id_USERS),
   FOREIGN KEY(Id_USERS) REFERENCES USERS(Id_USERS)
);

CREATE TABLE PREMISES(
   Id_PREMISES INT AUTO_INCREMENT,
   name VARCHAR(100)  NOT NULL,
   streetNumber SMALLINT,
   streetName VARCHAR(100) ,
   city VARCHAR(100) ,
   country VARCHAR(50) ,
   PRIMARY KEY(Id_PREMISES),
   UNIQUE(name)
);

CREATE TABLE COOKING_SPACES(
   Id_COOKING_SPACES INT AUTO_INCREMENT,
   name VARCHAR(50) ,
   size SMALLINT,
   isAvailable BOOLEAN DEFAULT FALSE,
   PricePerHour DECIMAL(19,4),
   Id_PREMISES INT NOT NULL,
   PRIMARY KEY(Id_COOKING_SPACES),
   FOREIGN KEY(Id_PREMISES) REFERENCES PREMISES(Id_PREMISES)
);

CREATE TABLE COOKING_ITEMS(
   Id_COOKING_ITEMS INT AUTO_INCREMENT,
   name VARCHAR(100) ,
   status VARCHAR(50) ,
   Id_COOKING_SPACES INT NOT NULL,
   PRIMARY KEY(Id_COOKING_ITEMS),
   FOREIGN KEY(Id_COOKING_SPACES) REFERENCES COOKING_SPACES(Id_COOKING_SPACES)
);

CREATE TABLE SUBSCRIPTIONS(
   Id_SUBSCRIPTIONS INT AUTO_INCREMENT,
   name VARCHAR(50) ,
   price DECIMAL(19,4),
   max_lesson_access INT,
   picture VARCHAR(50) ,
   description TEXT,
   PRIMARY KEY(Id_SUBSCRIPTIONS),
   UNIQUE(name)
);

CREATE TABLE CONVERSATIONS(
   Id_CONVERSATIONS INT AUTO_INCREMENT,
   isClosed BOOLEAN,
   Id_USERS2 INT,
   Id_USERS1 INT NOT NULL,
   PRIMARY KEY(Id_CONVERSATIONS)
);

CREATE TABLE MESSAGES(
   Id_MESSAGES INT AUTO_INCREMENT,
   content TEXT,
   createdAt DATETIME NOT NULL,
   idSender INT NOT NULL,
   Id_CONVERSATIONS INT NOT NULL,
   PRIMARY KEY(Id_MESSAGES),
   FOREIGN KEY(Id_CONVERSATIONS) REFERENCES CONVERSATIONS(Id_CONVERSATIONS)
);

CREATE TABLE SHOP_ITEMS(
   Id_SHOP_ITEMS INT AUTO_INCREMENT,
   name VARCHAR(100) ,
   description TEXT,
   price DECIMAL(19,4),
   stock BIGINT,
   reward VARCHAR(50) ,
   picture VARCHAR(255)  NOT NULL DEFAULT 'default.png',
   PRIMARY KEY(Id_SHOP_ITEMS)
);

CREATE TABLE FOODS(
   Id_FOODS INT AUTO_INCREMENT,
   name VARCHAR(100) ,
   description TEXT,
   price DECIMAL(19,4),
   picture VARCHAR(255)  NOT NULL DEFAULT 'default.png',
   PRIMARY KEY(Id_FOODS)
);

CREATE TABLE INGREDIENTS(
   Id_INGREDIENTS INT AUTO_INCREMENT,
   name VARCHAR(100)  NOT NULL,
   Allergen VARCHAR(50) ,
   Id_COOKING_SPACES INT NOT NULL,
   PRIMARY KEY(Id_INGREDIENTS),
   FOREIGN KEY(Id_COOKING_SPACES) REFERENCES COOKING_SPACES(Id_COOKING_SPACES)
);

CREATE TABLE GROUPS(
   Id_GROUPS INT AUTO_INCREMENT,
   name VARCHAR(100)  NOT NULL,
   PRIMARY KEY(Id_GROUPS),
   UNIQUE(name)
);

CREATE TABLE BILLS(
   Id_BILLS INT AUTO_INCREMENT,
   name VARCHAR(255)  NOT NULL,
   type VARCHAR(50)  NOT NULL,
   createdAt DATETIME NOT NULL DEFAULT NOW(),
   Id_USERS INT NOT NULL,
   PRIMARY KEY(Id_BILLS),
   UNIQUE(name),
   FOREIGN KEY(Id_USERS) REFERENCES USERS(Id_USERS)
);

CREATE TABLE LESSONS_GROUPS(
   Id_LESSONS_GROUPS INT AUTO_INCREMENT,
   name VARCHAR(100)  NOT NULL,
   PRIMARY KEY(Id_LESSONS_GROUPS),
   UNIQUE(name)
);

CREATE TABLE CONTRACTOR_TYPES(
   Id_CONTRACTOR_TYPES INT AUTO_INCREMENT,
   name VARCHAR(50)  NOT NULL,
   PRIMARY KEY(Id_CONTRACTOR_TYPES),
   UNIQUE(name)
);

CREATE TABLE CONTRACTORS(
   Id_CONTRACTORS INT AUTO_INCREMENT,
   presentation TEXT,
   contractStart DATETIME NOT NULL,
   contractEnd DATETIME NOT NULL,
   Id_CONTRACTOR_TYPES INT NOT NULL,
   Id_USERS INT NOT NULL,
   PRIMARY KEY(Id_CONTRACTORS),
   UNIQUE(Id_USERS),
   FOREIGN KEY(Id_CONTRACTOR_TYPES) REFERENCES CONTRACTOR_TYPES(Id_CONTRACTOR_TYPES),
   FOREIGN KEY(Id_USERS) REFERENCES USERS(Id_USERS)
);

CREATE TABLE EVENTS(
   Id_EVENTS INT AUTO_INCREMENT,
   name VARCHAR(100)  NOT NULL,
   type VARCHAR(50) ,
   endTime DATETIME,
   isClosed BOOLEAN DEFAULT FALSE,
   startTime DATETIME,
   isInternal BOOLEAN,
   isPrivate BOOLEAN,
   group_display_order INT,
   defaultPicture VARCHAR(255)  DEFAULT 'default.jpg',
   Id_GROUPS INT,
   PRIMARY KEY(Id_EVENTS),
   UNIQUE(name),
   FOREIGN KEY(Id_GROUPS) REFERENCES GROUPS(Id_GROUPS)
);

CREATE TABLE COMMENTS(
   Id_COMMENTS INT AUTO_INCREMENT,
   grade DECIMAL(2,1)   NOT NULL,
   content TEXT,
   picture VARCHAR(255)  DEFAULT 'default.png',
   Id_CLIENTS INT NOT NULL,
   Id_EVENTS INT NOT NULL,
   PRIMARY KEY(Id_COMMENTS),
   FOREIGN KEY(Id_CLIENTS) REFERENCES CLIENTS(Id_CLIENTS),
   FOREIGN KEY(Id_EVENTS) REFERENCES EVENTS(Id_EVENTS)
);

CREATE TABLE ORDERS(
   Id_ORDERS INT AUTO_INCREMENT,
   status VARCHAR(20) ,
   price DECIMAL(19,4) NOT NULL,
   deliveryAddress VARCHAR(50) ,
   Id_CONTRACTORS INT NOT NULL,
   Id_CONTRACTORS_1 INT NOT NULL,
   Id_CLIENTS INT NOT NULL,
   PRIMARY KEY(Id_ORDERS),
   FOREIGN KEY(Id_CONTRACTORS) REFERENCES CONTRACTORS(Id_CONTRACTORS),
   FOREIGN KEY(Id_CONTRACTORS_1) REFERENCES CONTRACTORS(Id_CONTRACTORS),
   FOREIGN KEY(Id_CLIENTS) REFERENCES CLIENTS(Id_CLIENTS)
);

CREATE TABLE LESSONS(
   Id_LESSONS INT AUTO_INCREMENT,
   name VARCHAR(50) ,
   content TEXT,
   description VARCHAR(50) ,
   difficulty TINYINT,
   group_display_order INT,
   Id_LESSONS_GROUPS INT,
   PRIMARY KEY(Id_LESSONS),
   FOREIGN KEY(Id_LESSONS_GROUPS) REFERENCES LESSONS_GROUPS(Id_LESSONS_GROUPS)
);

CREATE TABLE IS_SUBSCRIBED(
   Id_CLIENTS INT,
   Id_SUBSCRIPTIONS INT,
   endTime DATETIME NOT NULL,
   PRIMARY KEY(Id_CLIENTS, Id_SUBSCRIPTIONS),
   FOREIGN KEY(Id_CLIENTS) REFERENCES CLIENTS(Id_CLIENTS),
   FOREIGN KEY(Id_SUBSCRIPTIONS) REFERENCES SUBSCRIPTIONS(Id_SUBSCRIPTIONS)
);

CREATE TABLE PARTICIPATES(
   Id_CLIENTS INT,
   Id_EVENTS INT,
   isPresent BOOLEAN DEFAULT FALSE,
   PRIMARY KEY(Id_CLIENTS, Id_EVENTS),
   FOREIGN KEY(Id_CLIENTS) REFERENCES CLIENTS(Id_CLIENTS),
   FOREIGN KEY(Id_EVENTS) REFERENCES EVENTS(Id_EVENTS)
);

CREATE TABLE ANIMATES(
   Id_CONTRACTORS INT,
   Id_EVENTS INT,
   PRIMARY KEY(Id_CONTRACTORS, Id_EVENTS),
   FOREIGN KEY(Id_CONTRACTORS) REFERENCES CONTRACTORS(Id_CONTRACTORS),
   FOREIGN KEY(Id_EVENTS) REFERENCES EVENTS(Id_EVENTS)
);

CREATE TABLE ORGANIZES(
   Id_MANAGERS INT,
   Id_EVENTS INT,
   PRIMARY KEY(Id_MANAGERS, Id_EVENTS),
   FOREIGN KEY(Id_MANAGERS) REFERENCES MANAGERS(Id_MANAGERS),
   FOREIGN KEY(Id_EVENTS) REFERENCES EVENTS(Id_EVENTS)
);

CREATE TABLE IS_HOSTED(
   Id_COOKING_SPACES INT,
   Id_EVENTS INT,
   PRIMARY KEY(Id_COOKING_SPACES, Id_EVENTS),
   FOREIGN KEY(Id_COOKING_SPACES) REFERENCES COOKING_SPACES(Id_COOKING_SPACES),
   FOREIGN KEY(Id_EVENTS) REFERENCES EVENTS(Id_EVENTS)
);

CREATE TABLE TALKS(
   Id_USERS INT,
   Id_CONVERSATIONS INT,
   PRIMARY KEY(Id_USERS, Id_CONVERSATIONS),
   FOREIGN KEY(Id_USERS) REFERENCES USERS(Id_USERS),
   FOREIGN KEY(Id_CONVERSATIONS) REFERENCES CONVERSATIONS(Id_CONVERSATIONS)
);

CREATE TABLE TEACHES(
   Id_CONTRACTORS INT,
   Id_LESSONS INT,
   PRIMARY KEY(Id_CONTRACTORS, Id_LESSONS),
   FOREIGN KEY(Id_CONTRACTORS) REFERENCES CONTRACTORS(Id_CONTRACTORS),
   FOREIGN KEY(Id_LESSONS) REFERENCES LESSONS(Id_LESSONS)
);

CREATE TABLE CONTAINS_ITEM(
   Id_ORDERS INT,
   Id_SHOP_ITEMS INT,
   PRIMARY KEY(Id_ORDERS, Id_SHOP_ITEMS),
   FOREIGN KEY(Id_ORDERS) REFERENCES ORDERS(Id_ORDERS),
   FOREIGN KEY(Id_SHOP_ITEMS) REFERENCES SHOP_ITEMS(Id_SHOP_ITEMS)
);

CREATE TABLE CONTAINS_FOOD(
   Id_ORDERS INT,
   Id_FOODS INT,
   PRIMARY KEY(Id_ORDERS, Id_FOODS),
   FOREIGN KEY(Id_ORDERS) REFERENCES ORDERS(Id_ORDERS),
   FOREIGN KEY(Id_FOODS) REFERENCES FOODS(Id_FOODS)
);

CREATE TABLE WATCHES(
   Id_CLIENTS INT,
   Id_LESSONS INT,
   PRIMARY KEY(Id_CLIENTS, Id_LESSONS),
   FOREIGN KEY(Id_CLIENTS) REFERENCES CLIENTS(Id_CLIENTS),
   FOREIGN KEY(Id_LESSONS) REFERENCES LESSONS(Id_LESSONS)
);

CREATE TABLE BOOKS(
   Id_CLIENTS INT,
   Id_COOKING_SPACES INT,
   startTime DATETIME NOT NULL,
   endTime DATETIME NOT NULL,
   PRIMARY KEY(Id_CLIENTS, Id_COOKING_SPACES),
   FOREIGN KEY(Id_CLIENTS) REFERENCES CLIENTS(Id_CLIENTS),
   FOREIGN KEY(Id_COOKING_SPACES) REFERENCES COOKING_SPACES(Id_COOKING_SPACES)
);
