# DynamicStockManager
Case Study for Perwatch

    // Form related APIs
	Post[/api/v1/form/create]
	GET[/api/v1/form]
	GET[/api/v1/form/:_id]
	PUT[/api/v1/form/:_id]
	DELETE[/api/v1/form/:_id]

	// Field related APIs
	POST[/api/v1/form/:_id/field]
	GET[/api/v1/form/:_id/field]
	GET[/api/v1/form/:_id/field/:field_id]
	PUT[/api/v1/form/:_id/field/:field_id]
    DELETE[/api/v1/form/:_id/field/:field_id]

	// Stock related APIs

	POST[/api/v1/form/:_id/stock]
	GET[/api/v1/form/:_id/stock]
	GET[/api/v1/form/:_id/stock/:stock_id]
	PUT[/api/v1/form/:_id/stock/:stock_id]
	DELETE[/api/v1/form/:_id/stock/:stock_id]