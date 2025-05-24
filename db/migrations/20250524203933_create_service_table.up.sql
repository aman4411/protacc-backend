-- Create service status enum type
CREATE TYPE service_status AS ENUM ('active', 'inactive');

-- Create service categories table
CREATE TABLE service_categories (
                                    id SERIAL PRIMARY KEY,
                                    name VARCHAR(100) NOT NULL,
                                    slug VARCHAR(100) NOT NULL UNIQUE,
                                    description TEXT,
                                    icon VARCHAR(50),
                                    status service_status DEFAULT 'active',
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create services table
CREATE TABLE services (
                          id SERIAL PRIMARY KEY,
                          category_id INTEGER REFERENCES service_categories(id),
                          name VARCHAR(200) NOT NULL,
                          slug VARCHAR(200) NOT NULL UNIQUE,
                          description TEXT,
                          short_description VARCHAR(500),
                          features TEXT[],
                          requirements TEXT[],
                          price DECIMAL(10,2) NOT NULL,
                          booking_amount DECIMAL(10,2) DEFAULT 99.00,
                          estimated_delivery_days INTEGER,
                          icon VARCHAR(50),
                          status service_status DEFAULT 'active',
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create cart items table
CREATE TABLE cart_items (
                            id SERIAL PRIMARY KEY,
                            user_id UUID REFERENCES users(id),
                            service_id INTEGER REFERENCES services(id),
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            UNIQUE(user_id, service_id)
);

-- Create order status enum type
CREATE TYPE order_status AS ENUM (
    'pending_payment',
    'payment_received',
    'processing',
    'documents_required',
    'documents_received',
    'in_progress',
    'completed',
    'cancelled'
    );

-- Create orders table
CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        user_id UUID REFERENCES users(id),
                        service_id INTEGER REFERENCES services(id),
                        order_number VARCHAR(50) UNIQUE NOT NULL,
                        total_amount DECIMAL(10,2) NOT NULL,
                        booking_amount DECIMAL(10,2) NOT NULL,
                        remaining_amount DECIMAL(10,2) NOT NULL,
                        status order_status DEFAULT 'pending_payment',
                        payment_status BOOLEAN DEFAULT false,
                        notes TEXT,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create order status history table
CREATE TABLE order_status_history (
                                      id SERIAL PRIMARY KEY,
                                      order_id INTEGER REFERENCES orders(id),
                                      status order_status NOT NULL,
                                      notes TEXT,
                                      created_by UUID REFERENCES users(id),
                                      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_services_category_id ON services(category_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_service_id ON orders(service_id);
CREATE INDEX IF NOT EXISTS idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX IF NOT EXISTS idx_order_status_history_created_by ON order_status_history(created_by);