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

	       i32.const 40
	       local.set 0

	       i32.const 50
	       i32.const 60
	       i32.add

	       local.set 1
	       )
	local.get $l1
	local.get $l2
	i32.add
	i32.add
	)
  (export "main" (func $main)))
