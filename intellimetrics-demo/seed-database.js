db = db.getSiblingDB('intellimetrics');

db.companies.insertMany([
  {
    company_id: "COMP-10001",
    company_name: "TechStartup Inc",
    revenue: 1247000,
    users: 2847,
    growth_rate: 0.23,
    industry: "SaaS",
    created_at: new Date()
  },
  {
    company_id: "COMP-10002",
    company_name: "DataVision Corp",
    revenue: 3420000,
    users: 5120,
    growth_rate: 0.18,
    industry: "Analytics",
    created_at: new Date()
  },
  {
    company_id: "COMP-10003",
    company_name: "CloudScale Solutions",
    revenue: 890000,
    users: 1523,
    growth_rate: 0.45,
    industry: "Cloud Infrastructure",
    created_at: new Date()
  },
  {
    company_id: "COMP-10004",
    company_name: "AI Innovations Ltd",
    revenue: 2150000,
    users: 3890,
    growth_rate: 0.31,
    industry: "Artificial Intelligence",
    created_at: new Date()
  },
  {
    company_id: "COMP-10005",
    company_name: "SecureNet Systems",
    revenue: 1680000,
    users: 2340,
    growth_rate: 0.15,
    industry: "Cybersecurity",
    created_at: new Date()
  },
  {
    company_id: "COMP-10006",
    company_name: "FinTech Global",
    revenue: 4570000,
    users: 8920,
    growth_rate: 0.28,
    industry: "Financial Technology",
    created_at: new Date()
  },
  {
    company_id: "COMP-10007",
    company_name: "HealthTech Partners",
    revenue: 2980000,
    users: 4560,
    growth_rate: 0.22,
    industry: "Healthcare",
    created_at: new Date()
  },
  {
    company_id: "COMP-10008",
    company_name: "EduPlatform Co",
    revenue: 1120000,
    users: 6780,
    growth_rate: 0.35,
    industry: "EdTech",
    created_at: new Date()
  },
  {
    company_id: "COMP-10009",
    company_name: "RetailMetrics Inc",
    revenue: 3240000,
    users: 5450,
    growth_rate: 0.19,
    industry: "E-commerce",
    created_at: new Date()
  },
  {
    company_id: "COMP-10010",
    company_name: "GreenEnergy Analytics",
    revenue: 1890000,
    users: 2120,
    growth_rate: 0.41,
    industry: "Renewable Energy",
    created_at: new Date()
  }
]);

print("Database seeded with 10 companies");