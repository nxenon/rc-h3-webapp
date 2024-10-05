# Database Creation

```mysql
alter table CARTS
    drop foreign key FK_CARTS_RELATIONS_USERS;

alter table CART_PRODUCTS
    drop foreign key FK_CART_PRO_RELATIONS_PRODUCTS;

alter table CART_PRODUCTS
    drop foreign key FK_CART_PRO_RELATIONS_CARTS;

alter table RATE_LIMITS
    drop foreign key FK_RATE_LIM_RELATIONS_USERS;


alter table CARTS
    drop foreign key FK_CARTS_RELATIONS_USERS;

drop table if exists CARTS;


alter table CART_PRODUCTS
    drop foreign key FK_CART_PRO_RELATIONS_PRODUCTS;

alter table CART_PRODUCTS
    drop foreign key FK_CART_PRO_RELATIONS_CARTS;

drop table if exists CART_PRODUCTS;

drop table if exists COUPONS;

drop table if exists PRODUCTS;


alter table RATE_LIMITS
    drop foreign key FK_RATE_LIM_RELATIONS_USERS;

drop table if exists RATE_LIMITS;

drop table if exists USERS;

/*==============================================================*/
/* Table: CARTS                                                 */
/*==============================================================*/
create table CARTS
(
    CART_ID              int not null auto_increment,
    USER_ID              int,
    CART_OVERAL_PRICE    int not null,
    primary key (CART_ID)
);

/*==============================================================*/
/* Table: CART_PRODUCTS                                         */
/*==============================================================*/
create table CART_PRODUCTS
(
    PRODUCT_IN_CART_ID   int not null auto_increment,
    PRODUCT_ID           int,
    CART_ID              int,
    primary key (PRODUCT_IN_CART_ID)
);

/*==============================================================*/
/* Table: COUPONS                                               */
/*==============================================================*/
create table COUPONS
(
    COUPON_ID            int not null auto_increment,
    IS_VALID             bool not null,
    DISCOUNT_PERCENT     int not null,
    COUPON_VALUE         varchar(20) not null,
    primary key (COUPON_ID)
);

/*==============================================================*/
/* Table: PRODUCTS                                              */
/*==============================================================*/
create table PRODUCTS
(
    PRODUCT_ID           int not null auto_increment,
    PRODUCT_NAME         varchar(255) not null,
    PRODUCT_PRICE        int not null,
    IMAGE_URL            varchar(500) not null,
    primary key (PRODUCT_ID)
);

/*==============================================================*/
/* Table: RATE_LIMITS                                           */
/*==============================================================*/
create table RATE_LIMITS
(
    RATE_LIMIT_ID        int not null auto_increment,
    USER_ID              int,
    IP_ADDRESS           varchar(255),
    ATTEMPTS_NUMBER      int,
    primary key (RATE_LIMIT_ID)
);

/*==============================================================*/
/* Table: USERS                                                 */
/*==============================================================*/
create table USERS
(
    USER_ID              int not null auto_increment,
    USERNAME             varchar(255) not null,
    PASSWORD_HASH        varchar(255) not null,
    BALANCE              int not null,
    primary key (USER_ID)
);

alter table CARTS add constraint FK_CARTS_RELATIONS_USERS foreign key (USER_ID)
    references USERS (USER_ID) on delete restrict on update restrict;

alter table CART_PRODUCTS add constraint FK_CART_PRO_RELATIONS_PRODUCTS foreign key (PRODUCT_ID)
    references PRODUCTS (PRODUCT_ID) on delete restrict on update restrict;

alter table CART_PRODUCTS add constraint FK_CART_PRO_RELATIONS_CARTS foreign key (CART_ID)
    references CARTS (CART_ID) on delete restrict on update restrict;

alter table RATE_LIMITS add constraint FK_RATE_LIM_RELATIONS_USERS foreign key (USER_ID)
    references USERS (USER_ID) on delete restrict on update restrict;
```
