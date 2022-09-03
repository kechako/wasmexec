(module
  (func $main
	(result i32)

	(local $l1 i32)
	(local $l2 i32)

	i32.const 10
	local.set $l1

	i32.const 20
	local.set $l2

	i32.const 30
	(block $block1
	       (param i32)
	       (result i32)

	       i32.const 40
	       local.set 0

	       i32.const 50
	       local.set 1

	       i32.const 60
	       i32.add
	       )
	local.get $l1
	local.get $l2
	i32.add
	i32.add
	)
  (export "main" (func $main)))
