-- Clear existing services and categories to start fresh
TRUNCATE TABLE order_status_history, orders, cart_items, services, service_categories RESTART IDENTITY CASCADE;

-- Insert comprehensive service categories
INSERT INTO service_categories (name, slug, description, icon, status) VALUES
('Business Registration', 'business-registration', 'Complete business registration and startup services for various entity types', '/images/categories/business-registration.svg', 'active'),
('Government Registrations', 'government-registrations', 'Essential government registrations and certifications for businesses', '/images/categories/government-registrations.svg', 'active'),
('GST Services', 'gst-services', 'Comprehensive GST registration, filing, and compliance services', '/images/categories/gst-services.svg', 'active'),
('Income Tax Services', 'income-tax-services', 'Complete income tax filing, PAN, TAN, and compliance services', '/images/categories/income-tax-services.svg', 'active'),
('Legal Notice Handling', 'legal-notice-handling', 'Expert assistance for GST, Income Tax, and trademark notices', '/images/categories/legal-notices.svg', 'active'),
('Compliance Services', 'compliance-services', 'Ongoing compliance and return filing services for businesses', '/images/categories/compliance-services.svg', 'active'),
('Additional Services', 'additional-services', 'Consultancy, documentation, and other professional services', '/images/categories/additional-services.svg', 'active');

-- Insert startup/business registration services
INSERT INTO services (
    category_id, name, slug, description, short_description, features, requirements, 
    price, booking_amount, estimated_delivery_days, icon, status
) VALUES

-- Startup Services
(1, 'Proprietorship Registration', 'proprietorship-registration',
 'Register your sole proprietorship business with all necessary compliances and documentation support.',
 'Quick and easy proprietorship registration with expert guidance',
 ARRAY['Business name registration', 'Shop & Establishment license', 'GST registration assistance', 'PAN card assistance', 'Bank account opening guidance'],
 ARRAY['Owner ID proof', 'Address proof', 'Business address proof', 'Passport size photographs'],
 2999.00, 99.00, 7, 'user', 'active'),

(1, 'Partnership Registration', 'partnership-registration',
 'Complete partnership firm registration with partnership deed and all legal formalities.',
 'Professional partnership firm registration with legal documentation',
 ARRAY['Partnership deed drafting', 'Firm registration', 'PAN & TAN application', 'GST registration', 'Bank account opening assistance'],
 ARRAY['Partners ID proofs', 'Address proofs', 'Passport photographs', 'Business premises proof'],
 4999.00, 99.00, 10, 'users', 'active'),

(1, 'LLP Registration', 'llp-registration',
 'Limited Liability Partnership registration with complete documentation and compliance setup.',
 'Register your LLP with professional assistance and ongoing support',
 ARRAY['LLP incorporation', 'DPIN & DSC for partners', 'LLP agreement drafting', 'ROC registration', 'GST & other registrations'],
 ARRAY['Partners ID proofs', 'Address proofs', 'Digital signatures', 'Registered office proof'],
 7999.00, 99.00, 15, 'handshake', 'active'),

(1, 'One Person Company Registration', 'opc-registration',
 'OPC registration for single entrepreneurs with limited liability benefits.',
 'Start your OPC with complete legal protection and professional guidance',
 ARRAY['Company incorporation', 'Director DSC & DIN', 'MOA & AOA drafting', 'ROC registration', 'Compliance setup'],
 ARRAY['Director ID proof', 'Address proof', 'Passport photographs', 'Registered office proof'],
 8999.00, 99.00, 15, 'user-tie', 'active'),

(1, 'Private Limited Company Registration', 'private-limited-company',
 'Complete private limited company registration with all statutory compliances and expert guidance.',
 'Professional Pvt Ltd company registration with end-to-end support',
 ARRAY['Company name approval', 'DSC & DIN for directors', 'ROC incorporation', 'PAN, TAN & GST', 'Bank account assistance', 'Compliance calendar'],
 ARRAY['Directors ID proofs', 'Address proofs', 'Passport photographs', 'Digital signatures', 'Registered office documents'],
 9999.00, 99.00, 15, 'building', 'active'),

(1, 'Public Limited Company Registration', 'public-limited-company',
 'Public limited company incorporation for businesses planning to raise funds from public.',
 'Professional public company registration with IPO readiness guidance',
 ARRAY['Company incorporation', 'SEBI compliance setup', 'Share allotment', 'Statutory registers', 'Ongoing compliance'],
 ARRAY['Directors documents', 'Registered office proof', 'Capital structure details', 'Business plan'],
 24999.00, 999.00, 25, 'chart-line', 'active'),

(1, 'Section 8 Company Registration', 'section-8-ngo-registration',
 'Non-profit organization registration under Section 8 for charitable activities.',
 'Register your NGO/non-profit organization with tax exemptions',
 ARRAY['Section 8 license', 'Company incorporation', '12A & 80G registration', 'FCRA guidance', 'Compliance setup'],
 ARRAY['Promoters documents', 'Object clause', 'Office address proof', 'No-profit declaration'],
 12999.00, 299.00, 20, 'heart', 'active'),

(1, 'Trust Registration', 'trust-registration',
 'Trust registration for charitable, religious, or educational purposes with tax benefits.',
 'Complete trust registration with legal documentation and tax exemptions',
 ARRAY['Trust deed drafting', 'Trust registration', '12A & 80G registration', 'PAN application', 'Compliance guidance'],
 ARRAY['Trustees documents', 'Trust property details', 'Object clause', 'Address proof'],
 6999.00, 199.00, 12, 'shield-alt', 'active'),

(1, 'Producer Company Registration', 'producer-company-registration',
 'Producer company registration for agricultural and farming business collective.',
 'Register producer company for agricultural cooperatives and farming',
 ARRAY['Producer company incorporation', 'FPO registration', 'NABARD compliance', 'Subsidy guidance', 'Member management'],
 ARRAY['Producers list', 'Agricultural activity proof', 'Directors documents', 'Registered office proof'],
 14999.00, 399.00, 18, 'leaf', 'active'),

(1, 'Society Registration', 'society-registration',
 'Society registration for social, cultural, or educational activities with legal recognition.',
 'Professional society registration with legal compliance and tax benefits',
 ARRAY['Society registration', 'Memorandum & rules drafting', '12A & 80G registration', 'FCRA guidance', 'Annual compliance'],
 ARRAY['Members list', 'Office address proof', 'Object clause', 'Governing body details'],
 5999.00, 149.00, 10, 'users-cog', 'active'),

(1, 'Co-operative Society Registration', 'cooperative-society-registration',
 'Cooperative society registration for member-based business organizations.',
 'Register cooperative society with member benefits and legal structure',
 ARRAY['Society incorporation', 'Bye-laws drafting', 'Member registration', 'Government approvals', 'Compliance setup'],
 ARRAY['Members list', 'Office proof', 'Business plan', 'Governing body documents'],
 8999.00, 249.00, 15, 'handshake', 'active'),

-- Government Registration Services
(2, 'Startup India Registration', 'startup-india',
 'Startup India recognition for tax benefits, funding opportunities, and government support.',
 'Get official Startup India recognition with exclusive benefits and support',
 ARRAY['DPIIT recognition', 'Tax exemptions guidance', 'Funding assistance', 'IPR benefits', 'Regulatory support'],
 ARRAY['Company/LLP incorporation', 'Business plan', 'Innovation details', 'Financial projections'],
 4999.00, 99.00, 12, 'rocket', 'active'),

(2, 'Shop Act Registration', 'shop-act-registration',
 'Shop and Establishment Act registration for retail businesses and commercial establishments.',
 'Mandatory shop act license for all commercial establishments and shops',
 ARRAY['Shop act license', 'State compliance', 'Renewal assistance', 'Amendment support', 'Legal compliance'],
 ARRAY['Business address proof', 'Owner ID proof', 'NOC from landlord', 'Business activity details'],
 1999.00, 99.00, 5, 'store', 'active'),

(2, 'FSSAI Registration', 'fssai-registration',
 'Food Safety and Standards Authority registration for food business operators.',
 'FSSAI basic registration for small food businesses and startups',
 ARRAY['FSSAI basic registration', 'Food license guidance', 'Compliance support', 'Renewal assistance', 'Amendment help'],
 ARRAY['Business documents', 'Food category details', 'Premises proof', 'Owner ID proof'],
 2499.00, 99.00, 7, 'utensils', 'active'),

(2, 'FSSAI License', 'fssai-license',
 'FSSAI state or central license for medium to large scale food business operations.',
 'FSSAI state/central license for medium and large food businesses',
 ARRAY['FSSAI state/central license', 'Product approval', 'Lab testing guidance', 'Compliance support', 'Renewal assistance'],
 ARRAY['Detailed business plan', 'Premises layout', 'Equipment details', 'Product formulations', 'Quality manuals'],
 7999.00, 199.00, 15, 'certificate', 'active'),

(2, 'Import Export Code Registration', 'import-export-code',
 'IEC registration for businesses involved in import and export activities.',
 'Essential IEC code for international trade and import-export business',
 ARRAY['IEC code registration', 'DGFT portal setup', 'Export benefits guidance', 'Documentation support', 'Amendment assistance'],
 ARRAY['Company/firm documents', 'Bank certificate', 'PAN card', 'Address proof', 'Director/partner details'],
 3999.00, 99.00, 7, 'globe', 'active'),

(2, 'ICEGATE Registration', 'icegate-registration',
 'ICEGATE registration for customs clearance and trade facilitation.',
 'ICEGATE portal registration for seamless customs and trade operations',
 ARRAY['ICEGATE registration', 'Digital signature setup', 'Customs documentation', 'Trade facilitation', 'Portal training'],
 ARRAY['IEC code', 'Company documents', 'Authorized signatory details', 'Digital signature'],
 2999.00, 99.00, 5, 'ship', 'active'),

(2, 'LEI Code Registration', 'lei-registration',
 'Legal Entity Identifier code for financial transactions and regulatory compliance.',
 'LEI code for legal entities involved in financial market transactions',
 ARRAY['LEI code registration', 'Annual renewal', 'Regulatory compliance', 'Global recognition', 'Financial reporting'],
 ARRAY['Legal entity documents', 'Financial statements', 'Ownership structure', 'Management details'],
 4999.00, 149.00, 10, 'fingerprint', 'active'),

(2, 'ISO Certification', 'iso-registration',
 'ISO certification for quality management and international standards compliance.',
 'ISO certification for quality assurance and international business standards',
 ARRAY['ISO consultation', 'Documentation support', 'Audit assistance', 'Certification guidance', 'Maintenance support'],
 ARRAY['Business processes', 'Quality manual', 'Organization structure', 'Management commitment'],
 19999.00, 999.00, 45, 'award', 'active'),

(2, 'Trademark Registration', 'trademark-registration',
 'Trademark registration for brand protection and intellectual property rights.',
 'Protect your brand with official trademark registration and legal rights',
 ARRAY['Trademark search', 'Application filing', 'Examination response', 'Publication handling', 'Registration certificate'],
 ARRAY['Brand name/logo', 'Business documents', 'User affidavit', 'Power of attorney'],
 6999.00, 199.00, 60, 'registered', 'active'),

(2, 'Brand Name Registration', 'brand-name-registration',
 'Brand name registration and protection for business identity and marketing.',
 'Secure your brand name with legal protection and exclusive rights',
 ARRAY['Brand name search', 'Registration filing', 'Legal protection', 'Renewal assistance', 'Infringement support'],
 ARRAY['Brand details', 'Business documents', 'Usage proof', 'Authorization letter'],
 5999.00, 149.00, 45, 'tag', 'active'),

(2, 'Logo Registration', 'logo-registration',
 'Logo registration for visual brand protection and trademark rights.',
 'Protect your logo design with official registration and legal rights',
 ARRAY['Logo search', 'Design filing', 'Artistic work protection', 'Registration certificate', 'Renewal support'],
 ARRAY['Logo design', 'Business documents', 'Designer authorization', 'Usage evidence'],
 6999.00, 199.00, 50, 'palette', 'active'),

(2, 'ESI Registration', 'esi-registration',
 'Employee State Insurance registration for employee medical benefits and social security.',
 'ESI registration for employee welfare and mandatory social security compliance',
 ARRAY['ESI registration', 'Employee enrollment', 'Monthly returns', 'Compliance support', 'Benefit guidance'],
 ARRAY['Establishment documents', 'Employee list', 'Salary details', 'Bank account details'],
 2999.00, 99.00, 7, 'user-shield', 'active'),

(2, 'PF Registration', 'pf-registration',
 'Provident Fund registration for employee retirement benefits and statutory compliance.',
 'PF registration for employee retirement benefits and mandatory compliance',
 ARRAY['PF registration', 'Employee enrollment', 'Monthly contributions', 'Compliance support', 'Withdrawal assistance'],
 ARRAY['Establishment documents', 'Employee details', 'Salary structure', 'Bank account information'],
 2999.00, 99.00, 7, 'piggy-bank', 'active'),

(2, 'Udyam Registration', 'udyam-msme-registration',
 'Udyam registration for MSME benefits, subsidies, and government scheme eligibility.',
 'MSME Udyam registration for small business benefits and government support',
 ARRAY['Udyam certificate', 'MSME benefits', 'Subsidy guidance', 'Loan assistance', 'Tender benefits'],
 ARRAY['Business documents', 'Investment details', 'Turnover information', 'Bank statement'],
 1499.00, 99.00, 3, 'industry', 'active');

-- Insert GST Services
INSERT INTO services (
    category_id, name, slug, description, short_description, features, requirements, 
    price, booking_amount, estimated_delivery_days, icon, status
) VALUES

(3, 'GST Registration', 'gst-registration',
 'GST registration for businesses with turnover above threshold or voluntary registration.',
 'Complete GST registration with expert guidance and ongoing support',
 ARRAY['GST registration', 'GSTIN certificate', 'GST portal setup', 'Compliance calendar', 'Return filing guidance'],
 ARRAY['Business registration proof', 'PAN card', 'Address proof', 'Bank account details', 'Business photographs'],
 2999.00, 99.00, 7, 'file-invoice-dollar', 'active'),

(3, 'GST LUT Form', 'gst-lut-form',
 'Letter of Undertaking for GST-free exports without payment of integrated tax.',
 'GST LUT form for export businesses to avoid IGST payment',
 ARRAY['LUT application', 'Export documentation', 'Bank guarantee guidance', 'Compliance support', 'Renewal assistance'],
 ARRAY['GST registration', 'Export license', 'Bank details', 'CA certificate'],
 3999.00, 99.00, 10, 'file-export', 'active'),

(3, 'GST Amendment', 'gst-amendment',
 'GST registration amendment for changes in business details, address, or activities.',
 'Modify your GST registration details with proper documentation',
 ARRAY['Registration amendment', 'Supporting documents', 'Approval process', 'Updated certificate', 'Compliance update'],
 ARRAY['Amendment details', 'Supporting documents', 'Justification letter', 'Current GST certificate'],
 1999.00, 99.00, 15, 'edit', 'active'),

(3, 'GST Revocation', 'gst-revocation',
 'GST registration cancellation for businesses closing operations or below threshold.',
 'Cancel GST registration with proper closure procedures and compliance',
 ARRAY['Cancellation application', 'Final returns', 'Asset disposal', 'Liability clearance', 'Closure certificate'],
 ARRAY['Cancellation reason', 'Final accounts', 'Asset details', 'Liability clearance', 'Stock details'],
 2999.00, 99.00, 20, 'times-circle', 'active'),

(3, 'GST Number Transfer', 'gst-number-transfer',
 'GST registration transfer in case of business succession or ownership change.',
 'Transfer GST registration for business succession and ownership changes',
 ARRAY['Transfer application', 'Succession documents', 'New ownership proof', 'Liability transfer', 'Updated registration'],
 ARRAY['Succession documents', 'New owner details', 'Asset transfer deed', 'NOC from previous owner'],
 4999.00, 149.00, 25, 'exchange-alt', 'active'),

(3, 'GSTR-10 Filing', 'gstr-10',
 'GSTR-10 final return filing for GST registration cancellation or revocation.',
 'Final GST return filing for business closure and registration cancellation',
 ARRAY['Final return preparation', 'Asset disposal details', 'Liability settlement', 'Compliance closure', 'Certificate issuance'],
 ARRAY['Books of accounts', 'Asset details', 'Liability information', 'Stock valuation', 'Bank statements'],
 3999.00, 99.00, 15, 'file-alt', 'active'),

(3, 'GSTR-1 Filing', 'gstr-1-filing',
 'Monthly/Quarterly GSTR-1 return filing for outward supplies and invoices.',
 'Professional GSTR-1 return filing with invoice matching and compliance',
 ARRAY['Sales return filing', 'Invoice upload', 'Amendment handling', 'Error resolution', 'Compliance check'],
 ARRAY['Sales invoices', 'Credit/debit notes', 'Export documents', 'Previous returns'],
 1999.00, 99.00, 5, 'file-upload', 'active'),

(3, 'GSTR-3B Filing', 'gstr-3b-filing',
 'Monthly GSTR-3B return filing with tax liability calculation and payment.',
 'Accurate GSTR-3B filing with tax calculation and payment assistance',
 ARRAY['Monthly return filing', 'Tax calculation', 'Input credit optimization', 'Payment assistance', 'Late fee handling'],
 ARRAY['Purchase invoices', 'Sales data', 'Previous returns', 'Bank statements', 'Payment challans'],
 1499.00, 99.00, 3, 'calculator', 'active'),

(3, 'CMP-08 Filing', 'cmp-08-filing',
 'Composition scheme quarterly return filing for small businesses under GST.',
 'Simple CMP-08 return filing for composition scheme taxpayers',
 ARRAY['Quarterly return', 'Simple format', 'Turnover reporting', 'Tax calculation', 'Compliance check'],
 ARRAY['Sales data', 'Purchase details', 'Previous returns', 'Composition certificate'],
 999.00, 99.00, 2, 'file-signature', 'active'),

(3, 'GST Annual Return', 'gst-annual-return-r9',
 'Annual GST return filing with complete financial reconciliation and compliance.',
 'Comprehensive annual GST return with detailed reconciliation',
 ARRAY['Annual return preparation', 'Financial reconciliation', 'Audit trail', 'Compliance verification', 'Certificate generation'],
 ARRAY['All monthly returns', 'Financial statements', 'Audit reports', 'Bank statements', 'Purchase/sales registers'],
 7999.00, 299.00, 20, 'file-contract', 'active'),

(3, 'GST Audit', 'gst-audit-9c',
 'GST audit and reconciliation services for large taxpayers and compliance verification.',
 'Professional GST audit with detailed reconciliation and compliance report',
 ARRAY['Complete GST audit', 'Reconciliation report', 'Compliance verification', 'Discrepancy identification', 'Audit certificate'],
 ARRAY['All GST returns', 'Books of accounts', 'Financial statements', 'Supporting documents', 'Previous audit reports'],
 19999.00, 999.00, 30, 'search-dollar', 'active');

-- Insert Income Tax Services
INSERT INTO services (
    category_id, name, slug, description, short_description, features, requirements, 
    price, booking_amount, estimated_delivery_days, icon, status
) VALUES

(4, 'PAN Card Registration', 'pan-registration',
 'PAN card application for individuals and businesses for tax identification.',
 'Quick PAN card application with document verification and tracking',
 ARRAY['PAN application', 'Document verification', 'Application tracking', 'Dispatch notification', 'Correction assistance'],
 ARRAY['ID proof', 'Address proof', 'Date of birth proof', 'Passport size photographs'],
 599.00, 99.00, 10, 'id-card', 'active'),

(4, 'TAN Registration', 'tan-registration',
 'Tax Deduction Account Number registration for TDS deduction and compliance.',
 'TAN registration for businesses deducting TDS and tax compliance',
 ARRAY['TAN application', 'TDS compliance setup', 'Portal registration', 'Certificate download', 'Amendment support'],
 ARRAY['Business registration', 'PAN card', 'Address proof', 'Authorized signatory details'],
 1499.00, 99.00, 7, 'receipt', 'active'),

(4, 'ITR-1 Filing', 'itr-1-filing',
 'Income Tax Return filing for salaried individuals with simple income sources.',
 'Simple ITR-1 filing for salaried employees with basic income',
 ARRAY['ITR-1 preparation', 'Tax calculation', 'Refund processing', 'E-verification', 'Acknowledgment'],
 ARRAY['Form 16', 'Bank statements', 'Investment proofs', 'PAN card', 'Aadhaar card'],
 999.00, 99.00, 3, 'file-invoice', 'active'),

(4, 'ITR-2 Filing', 'itr-2-filing',
 'Income Tax Return filing for individuals with capital gains and multiple income sources.',
 'ITR-2 filing for individuals with house property and capital gains',
 ARRAY['ITR-2 preparation', 'Capital gains calculation', 'House property income', 'Tax optimization', 'E-filing'],
 ARRAY['Salary certificates', 'Capital gains documents', 'Property papers', 'Investment proofs', 'Bank statements'],
 1999.00, 99.00, 5, 'home', 'active'),

(4, 'ITR-3 Filing', 'itr-3-filing',
 'Income Tax Return filing for individuals with business or professional income.',
 'ITR-3 filing for proprietors and professionals with business income',
 ARRAY['ITR-3 preparation', 'Business income calculation', 'Professional income', 'Deduction optimization', 'Audit compliance'],
 ARRAY['Books of accounts', 'Profit & loss statement', 'Balance sheet', 'Business receipts', 'Expense vouchers'],
 2999.00, 99.00, 7, 'briefcase', 'active'),

(4, 'ITR-4 Filing', 'itr-4-filing',
 'Income Tax Return filing for presumptive taxation scheme under sections 44AD/44ADA.',
 'ITR-4 filing for small businesses under presumptive taxation scheme',
 ARRAY['ITR-4 preparation', 'Presumptive income calculation', 'Section 44AD/44ADA benefits', 'Simple compliance', 'Tax saving'],
 ARRAY['Business turnover details', 'Bank statements', 'Basic business records', 'PAN card'],
 1499.00, 99.00, 3, 'calculator', 'active'),

(4, 'ITR-5 Filing', 'itr-5-filing',
 'Income Tax Return filing for partnership firms, LLPs, and other entities.',
 'ITR-5 filing for firms, LLPs, and association of persons',
 ARRAY['ITR-5 preparation', 'Partnership income', 'LLP compliance', 'Partner allocation', 'Tax calculation'],
 ARRAY['Partnership deed', 'Profit & loss account', 'Balance sheet', 'Partners details', 'Books of accounts'],
 4999.00, 149.00, 10, 'handshake', 'active'),

(4, 'ITR-6 Filing', 'itr-6-filing',
 'Income Tax Return filing for companies under the Income Tax Act.',
 'ITR-6 filing for private and public limited companies',
 ARRAY['ITR-6 preparation', 'Company income calculation', 'Corporate tax compliance', 'Audit requirements', 'MAT calculation'],
 ARRAY['Audited financial statements', 'Tax audit report', 'Books of accounts', 'Board resolutions', 'Computation sheets'],
 9999.00, 299.00, 15, 'building', 'active'),

(4, 'ITR-7 Filing', 'itr-7-filing',
 'Income Tax Return filing for trusts, political parties, and exempt organizations.',
 'ITR-7 filing for trusts, institutions, and exempt entities',
 ARRAY['ITR-7 preparation', 'Exempt income handling', 'Trust compliance', 'Charitable activities', 'Exemption maintenance'],
 ARRAY['Trust deed', 'Activity reports', 'Financial statements', 'Exemption certificates', 'Donation receipts'],
 7999.00, 199.00, 12, 'shield-alt', 'active'),

(4, 'Form 15CA/15CB Filing', 'form-15ca-cb',
 'Form 15CA/15CB filing for foreign remittances and overseas transactions.',
 'Form 15CA/15CB filing for international payments and remittances',
 ARRAY['Form preparation', 'CA certification', 'Foreign exchange compliance', 'Documentation support', 'Bank coordination'],
 ARRAY['Payment details', 'Beneficiary information', 'Contract documents', 'Invoice details', 'PAN/TAN details'],
 2999.00, 99.00, 5, 'globe-americas', 'active'),

(4, 'TDS Return Filing', 'tds-return-filing',
 'TDS return filing for businesses deducting tax at source from payments.',
 'Professional TDS return filing with challan matching and compliance',
 ARRAY['TDS return preparation', 'Challan matching', 'Deductee details', 'Certificate generation', 'Correction assistance'],
 ARRAY['TDS challans', 'Deductee details', 'Payment records', 'Previous returns', 'Bank statements'],
 1999.00, 99.00, 5, 'receipt', 'active');

-- Insert Notice Handling Services
INSERT INTO services (
    category_id, name, slug, description, short_description, features, requirements, 
    price, booking_amount, estimated_delivery_days, icon, status
) VALUES

(5, 'GST Notice Handling', 'gst-notice-handling',
 'Professional assistance for GST notices, assessments, and compliance issues.',
 'Expert handling of GST notices with legal representation and resolution',
 ARRAY['Notice analysis', 'Legal representation', 'Response preparation', 'Hearing attendance', 'Appeal assistance'],
 ARRAY['GST notice copy', 'Business documents', 'Previous returns', 'Supporting evidence', 'Power of attorney'],
 7999.00, 299.00, 15, 'gavel', 'active'),

(5, 'Income Tax Notice Handling', 'income-tax-notice-handling',
 'Expert assistance for income tax notices, assessments, and proceedings.',
 'Professional income tax notice handling with legal support',
 ARRAY['Notice evaluation', 'Response drafting', 'Evidence compilation', 'Hearing representation', 'Appeal filing'],
 ARRAY['IT notice copy', 'Income tax returns', 'Supporting documents', 'Financial records', 'Authorization letter'],
 8999.00, 399.00, 20, 'balance-scale', 'active'),

(5, 'TDS Notice Handling', 'tds-notice-handling',
 'Specialized assistance for TDS notices, default proceedings, and compliance.',
 'TDS notice resolution with penalty minimization and compliance restoration',
 ARRAY['TDS notice analysis', 'Default resolution', 'Penalty negotiation', 'Compliance restoration', 'Future guidance'],
 ARRAY['TDS notice', 'TDS returns', 'Challan details', 'Deductee information', 'Payment records'],
 5999.00, 199.00, 12, 'hand-holding-usd', 'active'),

(5, 'Trademark Objection Handling', 'trademark-objection',
 'Professional response to trademark examination objections and oppositions.',
 'Expert trademark objection handling with legal arguments and evidence',
 ARRAY['Objection analysis', 'Legal response', 'Evidence compilation', 'Hearing representation', 'Registration completion'],
 ARRAY['Objection letter', 'Trademark application', 'Evidence of use', 'Legal authorization', 'Supporting documents'],
 9999.00, 499.00, 30, 'exclamation-triangle', 'active'),

(5, 'Brand Name Objection Handling', 'brand-name-objection',
 'Resolution of brand name examination objections and legal proceedings.',
 'Brand name objection resolution with trademark expertise',
 ARRAY['Objection review', 'Response preparation', 'Legal arguments', 'Evidence submission', 'Approval assistance'],
 ARRAY['Objection notice', 'Brand application', 'Usage evidence', 'Legal documents', 'Authorization'],
 7999.00, 299.00, 25, 'question-circle', 'active'),

(5, 'Logo Objection Handling', 'logo-objection',
 'Professional handling of logo design objections and trademark issues.',
 'Logo objection resolution with design and legal expertise',
 ARRAY['Design analysis', 'Objection response', 'Artistic arguments', 'Legal compliance', 'Registration support'],
 ARRAY['Objection letter', 'Logo application', 'Design evidence', 'Artistic documentation', 'Legal authorization'],
 8999.00, 399.00, 28, 'paint-brush', 'active');

-- Insert Compliance Services
INSERT INTO services (
    category_id, name, slug, description, short_description, features, requirements, 
    price, booking_amount, estimated_delivery_days, icon, status
) VALUES

(6, 'Company Annual Compliances', 'company-compliances',
 'Complete annual compliance services for private and public limited companies.',
 'Comprehensive company compliance including ROC filings and statutory requirements',
 ARRAY['Annual return filing', 'Financial statement filing', 'Board resolutions', 'Statutory registers', 'ROC compliance'],
 ARRAY['Financial statements', 'Board minutes', 'Shareholding details', 'Audit reports', 'Previous filings'],
 12999.00, 499.00, 25, 'clipboard-check', 'active'),

(6, 'LLP Annual Compliances', 'llp-compliances',
 'Annual compliance services for Limited Liability Partnerships including ROC filings.',
 'Complete LLP compliance with annual filings and statutory requirements',
 ARRAY['Annual return filing', 'Statement of accounts', 'Partner changes', 'Statutory compliance', 'ROC filings'],
 ARRAY['LLP accounts', 'Partner details', 'Previous returns', 'Audit reports', 'Partnership changes'],
 7999.00, 299.00, 20, 'handshake', 'active'),

(6, 'FSSAI Return Filing', 'fssai-return-filing',
 'FSSAI annual return filing for food business operators and license holders.',
 'FSSAI compliance with annual return filing and renewal assistance',
 ARRAY['Annual return filing', 'License renewal', 'Compliance verification', 'Amendment support', 'Penalty handling'],
 ARRAY['FSSAI license', 'Business details', 'Production data', 'Financial information', 'Previous returns'],
 2999.00, 99.00, 10, 'file-medical', 'active'),

(6, 'ESI Return Filing', 'esi-return-filing',
 'Monthly ESI return filing and compliance for employee social security.',
 'ESI monthly return filing with contribution calculation and compliance',
 ARRAY['Monthly ESI returns', 'Contribution calculation', 'Employee enrollment', 'Compliance monitoring', 'Penalty avoidance'],
 ARRAY['Employee details', 'Salary information', 'Previous returns', 'Challan details', 'Registration certificate'],
 1999.00, 99.00, 5, 'user-shield', 'active'),

(6, 'PF Return Filing', 'pf-return-filing',
 'Monthly PF return filing and compliance for employee provident fund.',
 'PF monthly return filing with contribution management and compliance',
 ARRAY['Monthly PF returns', 'Contribution tracking', 'Employee management', 'Compliance alerts', 'Withdrawal support'],
 ARRAY['Employee list', 'Salary details', 'Previous returns', 'Contribution challans', 'PF registration'],
 1999.00, 99.00, 5, 'piggy-bank', 'active');

-- Insert Additional Services
INSERT INTO services (
    category_id, name, slug, description, short_description, features, requirements, 
    price, booking_amount, estimated_delivery_days, icon, status
) VALUES

(7, 'Business Consultancy', 'consultancy',
 'Professional business consultancy for strategy, compliance, and growth planning.',
 'Expert business advice for strategy, compliance, and operational excellence',
 ARRAY['Business strategy', 'Compliance guidance', 'Financial planning', 'Growth consulting', 'Risk assessment'],
 ARRAY['Business details', 'Financial statements', 'Current challenges', 'Growth objectives', 'Market information'],
 4999.00, 199.00, 7, 'comments', 'active'),

(7, 'Project Report Preparation', 'project-reports',
 'Detailed project reports for loan applications, investors, and business planning.',
 'Professional project reports for funding, loans, and business proposals',
 ARRAY['Market research', 'Financial projections', 'Technical feasibility', 'Risk analysis', 'Executive summary'],
 ARRAY['Business concept', 'Market data', 'Financial requirements', 'Technical specifications', 'Team details'],
 14999.00, 999.00, 20, 'chart-bar', 'active'),

(7, 'CMA Data Preparation', 'cma-data',
 'Credit Monitoring Arrangement data preparation for loan applications and renewals.',
 'CMA data preparation for bank loan applications and credit facilities',
 ARRAY['CMA format preparation', 'Financial projections', 'Cash flow analysis', 'Ratio analysis', 'Bank presentation'],
 ARRAY['Financial statements', 'Loan requirements', 'Business projections', 'Bank formats', 'Previous CMA'],
 7999.00, 299.00, 10, 'chart-line', 'active'),

(7, 'Bookkeeping Services', 'bookkeeping',
 'Professional bookkeeping and accounting services for small and medium businesses.',
 'Complete bookkeeping solution with monthly reports and compliance',
 ARRAY['Daily bookkeeping', 'Monthly reports', 'Bank reconciliation', 'Invoice management', 'Expense tracking'],
 ARRAY['Business transactions', 'Bank statements', 'Invoices', 'Receipts', 'Previous records'],
 2999.00, 99.00, 30, 'book', 'active'),

(7, 'Partnership Deed Drafting', 'partnership-deed-drafting',
 'Professional partnership deed drafting with legal compliance and protection.',
 'Comprehensive partnership deed with legal protection and dispute resolution',
 ARRAY['Legal deed drafting', 'Partner rights definition', 'Profit sharing terms', 'Dispute resolution', 'Registration assistance'],
 ARRAY['Partner details', 'Business objectives', 'Capital contribution', 'Profit sharing ratio', 'Management structure'],
 3999.00, 149.00, 7, 'file-contract', 'active'),

(7, 'Rent Agreement Drafting', 'rent-agreement-drafting',
 'Legal rent agreement drafting for residential and commercial properties.',
 'Professional rent agreement with legal compliance and protection',
 ARRAY['Legal agreement drafting', 'Stamp duty guidance', 'Registration assistance', 'Clause customization', 'Renewal support'],
 ARRAY['Property details', 'Landlord information', 'Tenant details', 'Rent terms', 'Security deposit'],
 1999.00, 99.00, 3, 'home', 'active'),

(7, 'Digital Signature Certificate', 'digital-signature',
 'Class 2 and Class 3 digital signature certificates for legal and business use.',
 'Digital signature certificate for secure online transactions and filings',
 ARRAY['DSC application', 'Identity verification', 'Certificate installation', 'Usage training', 'Renewal assistance'],
 ARRAY['Identity proof', 'Address proof', 'PAN card', 'Photograph', 'Email/mobile verification'],
 1499.00, 99.00, 5, 'signature', 'active');
