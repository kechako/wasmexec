(module
  (func $main
	(result i32)
	(result i32)
	(result i32)
	(result i32)

	i32.const 5
	i32.const 20
	i32.add

	i32.const 5
	i32.const 2
	i32.sub

	i32.const 4
	i32.const 7
	i32.mul

	i32.const 5
	i32.const 2
	i32.div_s
	)
  (export "main" (func $main)))
