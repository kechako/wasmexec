(module
  (func $main
	(result i32)

	i32.const 20
	i32.const 10
	call $sub
	)
  (func $sub
	(param $p1 i32)
	(param $p2 i32)

	(result i32)

	(local $l1 i32)
	(local $l2 i32)

	local.get $p1
	local.get $p2
	i32.add
	local.set $l1

	local.get $p1
	local.get $p2
	i32.sub
	local.set $l2

	local.get $l1
	local.get $l2
	i32.mul

	return

	drop
	i32.const -1
	)
  (export "main" (func $main)))
