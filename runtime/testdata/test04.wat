(module
  (func $main
	(local $1 i32)
	(local $2 i32)
	(local $3 i32)
	(local $4 i32)

	(result i32)
	(result i32)
	(result i32)
	(result i32)
	(result i32)
	(result i32)
	(result i32)
	(result i32)
	(result i32)
	(result i32)
	(result i32)
	(result i32)

	i32.const 10
	local.set $1
	i32.const 20
	local.set $2
	i32.const 30
	local.set $3
	i32.const 40
	local.set $4

	local.get $4
	local.get $3
	local.get $2
	local.get $1

	i32.const 50
	local.tee $1
	i32.const 60
	local.tee $2
	i32.const 70
	local.tee $3
	i32.const 80
	local.tee $4

	local.get $4
	local.get $3
	local.get $2
	local.get $1
	)
  (export "main" (func $main)))
