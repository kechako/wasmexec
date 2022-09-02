(module
  (func $main
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
       	(result i32)
       	(result i32)
       	(result i32)
       	(result i32)
       	(result i32)
       	(result i32)

	i32.const 0
	i32.eqz

	i32.const 10
	i32.eqz

	i32.const 10
	i32.const 10
	i32.eq

	i32.const 10
	i32.const 20
	i32.eq

	i32.const 10
	i32.const 10
	i32.ne

	i32.const 10
	i32.const 20
	i32.ne

	i32.const 10
	i32.const 20
	i32.lt_s

	i32.const 10
	i32.const 10
	i32.lt_s

	i32.const 20
	i32.const 10
	i32.lt_s

	i32.const 10
	i32.const 20
	i32.gt_s

	i32.const 10
	i32.const 10
	i32.gt_s

	i32.const 20
	i32.const 10
	i32.gt_s

	i32.const 10
	i32.const 20
	i32.le_s

	i32.const 10
	i32.const 10
	i32.le_s

	i32.const 20
	i32.const 10
	i32.le_s

	i32.const 10
	i32.const 20
	i32.ge_s

	i32.const 10
	i32.const 10
	i32.ge_s

	i32.const 20
	i32.const 10
	i32.ge_s
	)
  (export "main" (func $main)))



