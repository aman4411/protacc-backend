-- Insert service categories
INSERT INTO service_categories (name, slug, description, icon, status) VALUES
                                                                           ('Business Registration', 'business-registration', 'Complete business registration services for various entity types', 'building', 'active'),
                                                                           ('Tax & Compliance', 'tax-compliance', 'Tax filing and compliance services for businesses', 'file-invoice', 'active'),
                                                                           ('Trademark & IP', 'trademark-ip', 'Intellectual property and trademark registration services', 'trademark', 'active'),
                                                                           ('Digital Services', 'digital-services', 'Digital certificates and online business services', 'laptop', 'active');

-- Insert services
INSERT INTO services (
    category_id,
    name,
    slug,
    description,
    short_description,
    features,
    requirements,
    price,
    booking_amount,
    estimated_delivery_days,
    icon,
    status
) VALUES
-- Business Registration Services
(1, 'Private Limited Company Registration', 'private-limited-company-registration',
 'Complete assistance in registering your Private Limited Company with all necessary compliances and documentation.',
 'Register your Private Limited Company with expert guidance and support',
 ARRAY['Company name approval', 'DSC and DIN for directors', 'ROC registration', 'GST registration', 'PAN & TAN', 'Bank account opening assistance'],
 ARRAY['Director ID proofs', 'Address proof', 'Passport size photographs', 'Digital signature', 'Proposed company address proof'],
 9999.00, 99.00, 15, 'building', 'active'),

(1, 'LLP Registration', 'llp-registration',
 'Limited Liability Partnership registration service with complete documentation and compliance support.',
 'Register your LLP with comprehensive support and guidance',
 ARRAY['Name approval', 'DSC for partners', 'DPIN for partners', 'LLP registration', 'PAN & TAN', 'Bank account assistance'],
 ARRAY['Partner ID proofs', 'Address proof', 'Photographs', 'Digital signature'],
 7999.00, 99.00, 12, 'handshake', 'active'),

-- Tax & Compliance Services
(2, 'GST Registration', 'gst-registration',
 'Complete GST registration service including application filing and certificate procurement.',
 'Get your GST registration done hassle-free',
 ARRAY['GST application filing', 'Document verification', 'Certificate procurement', 'Post-registration guidance'],
 ARRAY['PAN card', 'Address proof', 'Bank statement', 'Photograph'],
 2999.00, 99.00, 7, 'file-invoice', 'active'),

(2, 'Income Tax Return Filing', 'income-tax-return-filing',
 'Professional assistance in filing your income tax returns accurately and timely.',
 'Expert assistance for ITR filing',
 ARRAY['Income computation', 'Deductions planning', 'Tax calculation', 'Return filing', 'Post-filing support'],
 ARRAY['Form 16', 'Investment proofs', 'Bank statements', 'PAN card'],
 1999.00, 99.00, 5, 'calculator', 'active'),

-- Trademark & IP Services
(3, 'Trademark Registration', 'trademark-registration',
 'Complete trademark registration service including search, filing, and follow-up.',
 'Protect your brand with trademark registration',
 ARRAY['Trademark search', 'Application filing', 'Documentation', 'Follow-up support', 'Registration certificate'],
 ARRAY['Brand logo', 'Brand details', 'Business proof', 'ID proof'],
 6999.00, 99.00, 180, 'trademark', 'active'),

(3, 'Copyright Registration', 'copyright-registration',
 'Register your creative works with copyright protection.',
 'Protect your creative works with copyright registration',
 ARRAY['Work examination', 'Application filing', 'Documentation support', 'Registration certificate'],
 ARRAY['Copy of work', 'Creator details', 'ID proof'],
 4999.00, 99.00, 30, 'copyright', 'active'),

-- Digital Services
(4, 'Digital Signature Certificate', 'digital-signature-certificate',
 'Get your Digital Signature Certificate for business compliance.',
 'Obtain DSC for digital document signing',
 ARRAY['Class 2/3 DSC', 'USB token', 'Installation support', 'Usage guidance'],
 ARRAY['ID proof', 'Address proof', 'Photograph', 'Organization documents'],
 1499.00, 99.00, 3, 'signature', 'active'),

(4, 'MSME Registration', 'msme-registration',
 'Complete MSME/Udyam registration service for small businesses.',
 'Register your business under MSME/Udyam',
 ARRAY['Online application', 'Document verification', 'Certificate procurement', 'Benefits guidance'],
 ARRAY['Aadhaar card', 'PAN card', 'Bank details', 'Business details'],
 999.00, 99.00, 2, 'certificate', 'active');