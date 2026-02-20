DROP TABLE IF EXISTS additional_repair CASCADE;
DROP TABLE IF EXISTS parts_supply_service_order CASCADE;
DROP TABLE IF EXISTS service_service_order CASCADE;
DROP TABLE IF EXISTS payment CASCADE;
DROP TABLE IF EXISTS service_order CASCADE;
DROP TABLE IF EXISTS service_order_status CASCADE;
DROP TABLE IF EXISTS service CASCADE;
DROP TABLE IF EXISTS parts_supply CASCADE;
DROP TABLE IF EXISTS vehicle CASCADE;
DROP TABLE IF EXISTS customer CASCADE;
DROP TABLE IF EXISTS user CASCADE;
DROP TABLE IF EXISTS user_type CASCADE;
DROP TABLE if exists additional_repair_status CASCADE;

-- Base tables without dependencies

-- public.additional_repair_status definition

-- Drop table

-- DROP TABLE additional_repair_status;

CREATE TABLE additional_repair_status (
                                               id bigserial NOT NULL,
                                               description varchar(50) NOT NULL,
                                               CONSTRAINT additional_repair_status_pkey PRIMARY KEY (id)
);

-- public.user_type definition
CREATE TABLE user_type (
                                id bigserial NOT NULL,
                                "type" varchar(50) NOT NULL,
                                CONSTRAINT user_type_pkey PRIMARY KEY (id)
);

-- public.user definition
CREATE TABLE user (
                           id bigserial NOT NULL,
                           email varchar(100) NOT NULL,
                           "password" varchar(255) NOT NULL,
                           user_type_id int8 NOT NULL,
                           created_at timestamptz NULL,
                           updated_at timestamptz NULL,
                           deleted_at timestamptz NULL,
                           CONSTRAINT uni_user_email UNIQUE (email),
                           CONSTRAINT user_pkey PRIMARY KEY (id),
                           CONSTRAINT fk_user_type_users FOREIGN KEY (user_type_id) REFERENCES user_type(id)
);

-- public.service_order_status definition
CREATE TABLE service_order_status (
                                           id bigserial NOT NULL,
                                           description varchar(50) NOT NULL,
                                           CONSTRAINT service_order_status_pkey PRIMARY KEY (id)
);

-- public.parts_supply definition

-- Drop table

-- DROP TABLE parts_supply;

CREATE TABLE parts_supply (
                                   id bigserial NOT NULL,
                                   name varchar(100) NOT NULL,
                                   description text NULL,
                                   price numeric(10, 2) NOT NULL,
                                   CONSTRAINT parts_supply_pkey PRIMARY KEY (id)
);

-- public.customer definition with user reference

CREATE TABLE customer (
                             id BIGSERIAL PRIMARY KEY,
                             cpf_cnpj VARCHAR(20) UNIQUE NOT NULL,
                             fullname VARCHAR(255) NOT NULL,
                             phone_number VARCHAR(20),
                             user_id INTEGER NOT NULL,
                             FOREIGN KEY (user_id) REFERENCES user(id)
);

-- public.vehicle definition

CREATE TABLE vehicle (
                            id BIGSERIAL PRIMARY KEY,
                            plate VARCHAR(20) UNIQUE NOT NULL,
                            customer_id INTEGER NOT NULL,
                            model VARCHAR(100) NOT NULL,
                            brand VARCHAR(100) NOT NULL,
                            year INTEGER NOT NULL,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                            updated_at TIMESTAMP WITH TIME ZONE,
                            deleted_at TIMESTAMP WITH TIME ZONE,
                            FOREIGN KEY (customer_id) REFERENCES customer(id) ON DELETE CASCADE
);

-- public.service definition

CREATE TABLE service (
                              id BIGSERIAL PRIMARY KEY,
                              name VARCHAR(100) NOT NULL,
                              description TEXT,
                              price DECIMAL(10, 2) NOT NULL,
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                              updated_at TIMESTAMP WITH TIME ZONE,
                              deleted_at TIMESTAMP WITH TIME ZONE
);

-- public.service_order definition

CREATE TABLE service_order (
                                    id bigserial NOT NULL,
                                    customer_id int8 NOT NULL,
                                    vehicle_id int8 NOT NULL,
                                    os_status_id int8 NOT NULL,
                                    estimate numeric(10, 2) NULL,
                                    started_execution_date timestamptz NULL,
                                    final_execution_date timestamptz NULL,
                                    created_at timestamptz DEFAULT NOW(),
                                    updated_at timestamptz NULL,
                                    deleted_at timestamptz NULL,
                                    CONSTRAINT service_order_pkey PRIMARY KEY (id),
                                    CONSTRAINT fk_customer_service_order FOREIGN KEY (customer_id) REFERENCES customer(id),
                                    CONSTRAINT fk_vehicle_service_order FOREIGN KEY (vehicle_id) REFERENCES vehicle(id),
                                    CONSTRAINT fk_status_service_order FOREIGN KEY (os_status_id) REFERENCES service_order_status(id)
);

-- public.service_service_order definition (join table for services and service orders)
CREATE TABLE service_service_order (
                                            service_id int8 NOT NULL,
                                            service_order_id int8 NOT NULL,
                                            price numeric(10, 2) NOT NULL,
                                            CONSTRAINT service_service_order_pkey PRIMARY KEY (service_id, service_order_id),
                                            CONSTRAINT fk_service_service_order FOREIGN KEY (service_id) REFERENCES service(id),
                                            CONSTRAINT fk_service_order_services FOREIGN KEY (service_order_id) REFERENCES service_order(id)
);

-- public.parts_supply_service_order definition
CREATE TABLE parts_supply_service_order (
                                                 parts_supply_id int8 NOT NULL,
                                                 service_order_id int8 NOT NULL,
                                                 price numeric(10, 2) NOT NULL,
                                                 quantity int NOT NULL DEFAULT 1,
                                                 CONSTRAINT parts_supply_service_order_pkey PRIMARY KEY (parts_supply_id, service_order_id),
                                                 CONSTRAINT fk_parts_supply_service_order FOREIGN KEY (parts_supply_id) REFERENCES parts_supply(id),
                                                 CONSTRAINT fk_service_order_parts_supply FOREIGN KEY (service_order_id) REFERENCES service_order(id)
);

-- public.payment definition
CREATE TABLE payment (
                              id bigserial NOT NULL,
                              service_order_id int8 NOT NULL,
                              payment_date timestamptz NOT NULL,
                              amount numeric(10, 2) NOT NULL,
                              CONSTRAINT payment_pkey PRIMARY KEY (id),
                              CONSTRAINT uni_payment_service_order_id UNIQUE (service_order_id),
                              CONSTRAINT fk_service_order_payment FOREIGN KEY (service_order_id) REFERENCES service_order(id)
);

-- public.additional_repair definition

-- Drop table

-- DROP TABLE additional_repair;

CREATE TABLE additional_repair (
                                        id bigserial NOT NULL,
                                        service_order_id int8 NOT NULL,
                                        service_id int8 NOT NULL,
                                        parts_supply_id int8 NOT NULL,
                                        ar_status_id int8 NOT NULL,
                                        CONSTRAINT additional_repair_pkey PRIMARY KEY (id),
                                        CONSTRAINT fk_additional_repair_status_additional_repairs FOREIGN KEY (ar_status_id) REFERENCES additional_repair_status(id),
                                        CONSTRAINT fk_parts_supply_additional_repairs FOREIGN KEY (parts_supply_id) REFERENCES parts_supply(id),
                                        CONSTRAINT fk_service_additional_repairs FOREIGN KEY (service_id) REFERENCES service(id),
                                        CONSTRAINT fk_service_order_additional_repairs FOREIGN KEY (service_order_id) REFERENCES service_order(id)
);

-- Create indexes
CREATE INDEX idx_service_order_customer_id ON service_order(customer_id);
CREATE INDEX idx_service_order_vehicle_id ON service_order(vehicle_id);
CREATE INDEX idx_service_order_status_id ON service_order(os_status_id);
CREATE INDEX idx_service_order_deleted_at ON service_order(deleted_at);

-- Insert sample data
INSERT INTO user_type ("type")
VALUES
    ('ADMIN'),
    ('CUSTOMER'),
    ('MECHANIC');

INSERT INTO user (email, password, user_type_id, created_at)
VALUES
    ('admin@xpto.com', '$2a$10$YourHashedPasswordHere', 1, NOW()),
    ('joao@email.com', '$2a$10$YourHashedPasswordHere', 2, NOW()),
    ('joana@email.com', '$2a$10$YourHashedPasswordHere', 2, NOW());

INSERT INTO customer (cpf_cnpj, fullname, phone_number, user_id)
VALUES
    ('19134950869', 'João da Silva', '(11) 98765-4321', 1),
    ('93684508000137', 'Joana da Silva', '(11) 99775-4829', 2);

INSERT INTO vehicle (plate, customer_id, model, brand, year)
VALUES
    ('EXD6183', 1, 'Corsa', 'Chevrolet', 2005),
    ('AAA0X00', 2, 'Onix', 'Chevrolet', 2014);

INSERT INTO service (name, description, price)
VALUES
    ('Oil Change', 'Complete oil change service with filter replacement', 150.00),
    ('Brake Service', 'Brake pad replacement and system check', 300.00);

INSERT INTO service_order_status (description)
VALUES
    ('Pending'),
    ('In Progress'),
    ('Completed'),
    ('Cancelled');