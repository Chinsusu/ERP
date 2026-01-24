-- Seed data for QC checkpoints (default templates)
INSERT INTO qc_checkpoints (id, code, name, description, checkpoint_type, applies_to, test_items, is_active) VALUES
-- IQC - Incoming Quality Control for materials
('a0000001-0001-0001-0001-000000000001', 'IQC-RAW', 'Incoming Raw Material QC', 'Quality check for incoming raw materials', 'IQC', 'MATERIAL', 
 '[
   {"name": "Visual Inspection", "method": "Visual", "specification": "No discoloration, contamination", "type": "PASS_FAIL"},
   {"name": "Odor Check", "method": "Sensory", "specification": "Characteristic odor, no off-notes", "type": "PASS_FAIL"},
   {"name": "Certificate of Analysis", "method": "Document Review", "specification": "CoA matches specs", "type": "PASS_FAIL"},
   {"name": "pH Level", "method": "pH Meter", "specification": "Within spec range", "type": "NUMERIC", "min": 4.0, "max": 8.0},
   {"name": "Specific Gravity", "method": "Densitometer", "specification": "Within ±0.05", "type": "NUMERIC"}
 ]'::jsonb, true),

-- IQC for packaging materials
('a0000001-0001-0001-0001-000000000002', 'IQC-PKG', 'Incoming Packaging QC', 'Quality check for incoming packaging materials', 'IQC', 'MATERIAL',
 '[
   {"name": "Visual Inspection", "method": "Visual", "specification": "No defects, scratches, dents", "type": "PASS_FAIL"},
   {"name": "Dimension Check", "method": "Caliper", "specification": "Within tolerance", "type": "PASS_FAIL"},
   {"name": "Color Match", "method": "Light Box", "specification": "Matches approved sample", "type": "PASS_FAIL"},
   {"name": "Print Quality", "method": "Visual", "specification": "Clear, no smudging", "type": "PASS_FAIL"},
   {"name": "Sample Count", "method": "AQL Sampling", "specification": "AQL 1.0", "type": "NUMERIC"}
 ]'::jsonb, true),

-- IPQC - In-Process Quality Control
('a0000001-0001-0001-0001-000000000003', 'IPQC-MIX', 'In-Process Mixing QC', 'Quality check during mixing/blending process', 'IPQC', 'PRODUCT',
 '[
   {"name": "Temperature", "method": "Thermometer", "specification": "As per SOP", "type": "NUMERIC", "unit": "°C"},
   {"name": "Mixing Time", "method": "Timer", "specification": "As per SOP", "type": "NUMERIC", "unit": "min"},
   {"name": "Mixing Speed", "method": "Tachometer", "specification": "As per SOP", "type": "NUMERIC", "unit": "rpm"},
   {"name": "Appearance", "method": "Visual", "specification": "Homogeneous, no lumps", "type": "PASS_FAIL"},
   {"name": "pH Level", "method": "pH Meter", "specification": "Within target ±0.5", "type": "NUMERIC"}
 ]'::jsonb, true),

-- IPQC for filling
('a0000001-0001-0001-0001-000000000004', 'IPQC-FILL', 'In-Process Filling QC', 'Quality check during filling process', 'IPQC', 'PRODUCT',
 '[
   {"name": "Fill Weight/Volume", "method": "Scale", "specification": "Target ±2%", "type": "NUMERIC"},
   {"name": "Cap Torque", "method": "Torque Meter", "specification": "8-12 lb-in", "type": "NUMERIC", "min": 8, "max": 12},
   {"name": "Seal Integrity", "method": "Visual", "specification": "No leaks", "type": "PASS_FAIL"},
   {"name": "Label Placement", "method": "Visual", "specification": "Centered, no wrinkles", "type": "PASS_FAIL"}
 ]'::jsonb, true),

-- FQC - Final Quality Control
('a0000001-0001-0001-0001-000000000005', 'FQC-FINISH', 'Final Product QC', 'Final quality check before release', 'FQC', 'PRODUCT',
 '[
   {"name": "Appearance", "method": "Visual", "specification": "As per spec", "type": "PASS_FAIL"},
   {"name": "Color", "method": "Colorimeter", "specification": "Matches standard", "type": "PASS_FAIL"},
   {"name": "Odor", "method": "Sensory Panel", "specification": "Characteristic", "type": "PASS_FAIL"},
   {"name": "pH", "method": "pH Meter", "specification": "Per product spec", "type": "NUMERIC"},
   {"name": "Viscosity", "method": "Viscometer", "specification": "Per product spec", "type": "NUMERIC"},
   {"name": "Microbial Test", "method": "Lab Analysis", "specification": "TPC <100 CFU/g", "type": "PASS_FAIL"},
   {"name": "Stability Indicator", "method": "Accelerated Stability", "specification": "No separation", "type": "PASS_FAIL"},
   {"name": "Packaging Integrity", "method": "Visual", "specification": "No defects", "type": "PASS_FAIL"}
 ]'::jsonb, true);
