package entity

// DefaultCategories représente les catégories par défaut créées pour chaque nouvel utilisateur
var DefaultCategories = []Category{
	{Name: "Travail", Type: "task", Icon: "md:work", Color: "#FF6B6B"},
	{Name: "Études", Type: "task", Icon: "ion:book", Color: "#4ECDC4"},
	{Name: "Santé", Type: "task", Icon: "fa5:heartbeat", Color: "#FF7675"},
	{Name: "Sport", Type: "task", Icon: "mci:run", Color: "#74B9FF"},
	{Name: "Courses", Type: "task", Icon: "md:shopping-cart", Color: "#55A3FF"},
	{Name: "Maison", Type: "task", Icon: "ion:home", Color: "#A29BFE"},
	{Name: "Loisirs", Type: "task", Icon: "fa:gamepad", Color: "#FD79A8"},

	{Name: "Nourriture", Type: "expense", Icon: "md:restaurant", Color: "#FF6B6B"},
	{Name: "Transport", Type: "expense", Icon: "ion:car", Color: "#4ECDC4"},
	{Name: "Logement", Type: "expense", Icon: "fa5:house-user", Color: "#45B7D1"},
	{Name: "Santé", Type: "expense", Icon: "fa5:clinic-medical", Color: "#FF7675"},
	{Name: "Abonnements", Type: "expense", Icon: "mci:netflix", Color: "#E17055"},
	{Name: "Divertissement", Type: "expense", Icon: "md:movie", Color: "#FD79A8"},
	{Name: "Shopping", Type: "expense", Icon: "fa5:tshirt", Color: "#FDCB6E"},
	{Name: "Éducation", Type: "expense", Icon: "ion:school", Color: "#6C5CE7"},

	{Name: "Salaire", Type: "revenue", Icon: "fa5:money-check-alt", Color: "#00B894"},
	{Name: "Business", Type: "revenue", Icon: "ion:briefcase", Color: "#00CEC9"},
	{Name: "Investissements", Type: "revenue", Icon: "mci:chart-line", Color: "#74B9FF"},
	{Name: "Cadeaux", Type: "revenue", Icon: "ion:gift", Color: "#FD79A8"},
	{Name: "Remboursements", Type: "revenue", Icon: "fa5:hand-holding-usd", Color: "#55A3FF"},
}
