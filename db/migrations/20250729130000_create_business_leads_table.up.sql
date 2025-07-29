-- Create business_leads table
CREATE TABLE business_leads (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    company_name VARCHAR(255),
    business_type VARCHAR(100),
    services_interested TEXT[], -- Array of service types they're interested in
    budget_range VARCHAR(50),
    preferred_contact_method VARCHAR(20) CHECK (preferred_contact_method IN ('email', 'phone', 'both')) DEFAULT 'email',
    message TEXT,
    status VARCHAR(20) CHECK (status IN ('new', 'contacted', 'in_progress', 'qualified', 'converted', 'rejected', 'closed')) DEFAULT 'new',
    priority VARCHAR(20) CHECK (priority IN ('low', 'medium', 'high', 'urgent')) DEFAULT 'medium',
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    source VARCHAR(50) DEFAULT 'website',
    follow_up_date DATE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_business_leads_status ON business_leads(status);
CREATE INDEX idx_business_leads_priority ON business_leads(priority);
CREATE INDEX idx_business_leads_assigned_to ON business_leads(assigned_to);
CREATE INDEX idx_business_leads_created_at ON business_leads(created_at DESC);
CREATE INDEX idx_business_leads_email ON business_leads(email);
CREATE INDEX idx_business_leads_phone ON business_leads(phone);

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_business_leads_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_business_leads_updated_at_trigger
    BEFORE UPDATE ON business_leads
    FOR EACH ROW
    EXECUTE FUNCTION update_business_leads_updated_at(); 