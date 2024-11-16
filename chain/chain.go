package chain

// Runnable is an interface that defines a single method, Invoke, which takes a single input and returns a single output.
type Runnable[T1, T2 any] interface {
	Invoke(T1) (T2, error)
}

// fn is a type that wraps a function that takes a single input and returns a single output.
type fn[T1, T2 any] func(T1) T2

func (f fn[T1, T2]) Invoke(in T1) (out T2, err error) {
	return f(in), nil
}

// Func is a function that takes a function and returns a Runnable instance that wraps the function.
func Func[F ~func(T1) T2, T1, T2 any](f F) Runnable[T1, T2] {
	return fn[T1, T2](f)
}

// fn2 is a type that wraps a function that takes a single input and returns a single output and an error.
type fn2[T1, T2 any] func(T1) (T2, error)

func (f fn2[T1, T2]) Invoke(in T1) (out T2, err error) {
	return f(in)
}

// Func2 is a function that takes a function and returns a Runnable instance that wraps the function.
func Func2[F ~func(T1) (T2, error), T1, T2 any](f F) Runnable[T1, T2] {
	return fn2[T1, T2](f)
}

type chain2[T1, T2, T3 any] struct {
	r1 Runnable[T1, T2]
	r2 Runnable[T2, T3]
}

func (c chain2[T1, T2, T3]) Invoke(in T1) (out T3, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	return c.r2.Invoke(x)
}

// Chain2 takes 2 Runnable instances and returns a new Runnable instance that chains the two together.
func Chain2[R1 Runnable[T1, T2], R2 Runnable[T2, T3], T1, T2, T3 any](r1 R1, r2 R2) Runnable[T1, T3] {
	return chain2[T1, T2, T3]{
		r1: r1,
		r2: r2,
	}
}

type chain3[T1, T2, T3, T4 any] struct {
	r1 Runnable[T1, T2]
	r2 Runnable[T2, T3]
	r3 Runnable[T3, T4]
}

func (c chain3[T1, T2, T3, T4]) Invoke(in T1) (out T4, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	y, err := c.r2.Invoke(x)
	if err != nil {
		return
	}
	return c.r3.Invoke(y)
}

// Chain3 takes 3 Runnable instances and returns a new Runnable instance that chains the three together.
func Chain3[R1 Runnable[T1, T2], R2 Runnable[T2, T3], R3 Runnable[T3, T4], T1, T2, T3, T4 any](r1 R1, r2 R2, r3 R3) Runnable[T1, T4] {
	return chain3[T1, T2, T3, T4]{
		r1: r1,
		r2: r2,
		r3: r3,
	}
}

type chain4[T1, T2, T3, T4, T5 any] struct {
	r1 Runnable[T1, T2]
	r2 Runnable[T2, T3]
	r3 Runnable[T3, T4]
	r4 Runnable[T4, T5]
}

func (c chain4[T1, T2, T3, T4, T5]) Invoke(in T1) (out T5, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	y, err := c.r2.Invoke(x)
	if err != nil {
		return
	}
	z, err := c.r3.Invoke(y)
	if err != nil {
		return
	}
	return c.r4.Invoke(z)
}

// Chain4 takes 4 Runnable instances and returns a new Runnable instance that chains the four together.
func Chain4[R1 Runnable[T1, T2], R2 Runnable[T2, T3], R3 Runnable[T3, T4], R4 Runnable[T4, T5], T1, T2, T3, T4, T5 any](r1 R1, r2 R2, r3 R3, r4 R4) Runnable[T1, T5] {
	return chain4[T1, T2, T3, T4, T5]{
		r1: r1,
		r2: r2,
		r3: r3,
		r4: r4,
	}
}

type chain5[T1, T2, T3, T4, T5, T6 any] struct {
	r1 Runnable[T1, T2]
	r2 Runnable[T2, T3]
	r3 Runnable[T3, T4]
	r4 Runnable[T4, T5]
	r5 Runnable[T5, T6]
}

func (c chain5[T1, T2, T3, T4, T5, T6]) Invoke(in T1) (out T6, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	y, err := c.r2.Invoke(x)
	if err != nil {
		return
	}
	z, err := c.r3.Invoke(y)
	if err != nil {
		return
	}
	w, err := c.r4.Invoke(z)
	if err != nil {
		return
	}
	return c.r5.Invoke(w)
}

// Chain5 takes 5 Runnable instances and returns a new Runnable instance that chains the five together.
func Chain5[R1 Runnable[T1, T2], R2 Runnable[T2, T3], R3 Runnable[T3, T4], R4 Runnable[T4, T5], R5 Runnable[T5, T6], T1, T2, T3, T4, T5, T6 any](r1 R1, r2 R2, r3 R3, r4 R4, r5 R5) Runnable[T1, T6] {
	return chain5[T1, T2, T3, T4, T5, T6]{
		r1: r1,
		r2: r2,
		r3: r3,
		r4: r4,
		r5: r5,
	}
}

type chain6[T1, T2, T3, T4, T5, T6, T7 any] struct {
	r1 Runnable[T1, T2]
	r2 Runnable[T2, T3]
	r3 Runnable[T3, T4]
	r4 Runnable[T4, T5]
	r5 Runnable[T5, T6]
	r6 Runnable[T6, T7]
}

func (c chain6[T1, T2, T3, T4, T5, T6, T7]) Invoke(in T1) (out T7, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	y, err := c.r2.Invoke(x)
	if err != nil {
		return
	}
	z, err := c.r3.Invoke(y)
	if err != nil {
		return
	}
	w, err := c.r4.Invoke(z)
	if err != nil {
		return
	}
	u, err := c.r5.Invoke(w)
	if err != nil {
		return
	}
	return c.r6.Invoke(u)
}

// Chain6 takes 6 Runnable instances and returns a new Runnable instance that chains the six together.
func Chain6[R1 Runnable[T1, T2], R2 Runnable[T2, T3], R3 Runnable[T3, T4], R4 Runnable[T4, T5], R5 Runnable[T5, T6], R6 Runnable[T6, T7], T1, T2, T3, T4, T5, T6, T7 any](r1 R1, r2 R2, r3 R3, r4 R4, r5 R5, r6 R6) Runnable[T1, T7] {
	return chain6[T1, T2, T3, T4, T5, T6, T7]{
		r1: r1,
		r2: r2,
		r3: r3,
		r4: r4,
		r5: r5,
		r6: r6,
	}
}

type chain7[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	r1 Runnable[T1, T2]
	r2 Runnable[T2, T3]
	r3 Runnable[T3, T4]
	r4 Runnable[T4, T5]
	r5 Runnable[T5, T6]
	r6 Runnable[T6, T7]
	r7 Runnable[T7, T8]
}

func (c chain7[T1, T2, T3, T4, T5, T6, T7, T8]) Invoke(in T1) (out T8, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	y, err := c.r2.Invoke(x)
	if err != nil {
		return
	}
	z, err := c.r3.Invoke(y)
	if err != nil {
		return
	}
	w, err := c.r4.Invoke(z)
	if err != nil {
		return
	}
	u, err := c.r5.Invoke(w)
	if err != nil {
		return
	}
	v, err := c.r6.Invoke(u)
	if err != nil {
		return
	}
	return c.r7.Invoke(v)
}

// Chain7 takes 7 Runnable instances and returns a new Runnable instance that chains the seven together.
func Chain7[R1 Runnable[T1, T2], R2 Runnable[T2, T3], R3 Runnable[T3, T4], R4 Runnable[T4, T5], R5 Runnable[T5, T6], R6 Runnable[T6, T7], R7 Runnable[T7, T8], T1, T2, T3, T4, T5, T6, T7, T8 any](r1 R1, r2 R2, r3 R3, r4 R4, r5 R5, r6 R6, r7 R7) Runnable[T1, T8] {
	return chain7[T1, T2, T3, T4, T5, T6, T7, T8]{
		r1: r1,
		r2: r2,
		r3: r3,
		r4: r4,
		r5: r5,
		r6: r6,
		r7: r7,
	}
}

type chain8[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
	r1 Runnable[T1, T2]
	r2 Runnable[T2, T3]
	r3 Runnable[T3, T4]
	r4 Runnable[T4, T5]
	r5 Runnable[T5, T6]
	r6 Runnable[T6, T7]
	r7 Runnable[T7, T8]
	r8 Runnable[T8, T9]
}

func (c chain8[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Invoke(in T1) (out T9, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	y, err := c.r2.Invoke(x)
	if err != nil {
		return
	}
	z, err := c.r3.Invoke(y)
	if err != nil {
		return
	}
	w, err := c.r4.Invoke(z)
	if err != nil {
		return
	}
	u, err := c.r5.Invoke(w)
	if err != nil {
		return
	}
	v, err := c.r6.Invoke(u)
	if err != nil {
		return
	}
	t, err := c.r7.Invoke(v)
	if err != nil {
		return
	}
	return c.r8.Invoke(t)
}

// Chain8 takes 8 Runnable instances and returns a new Runnable instance that chains the eight together.
func Chain8[R1 Runnable[T1, T2], R2 Runnable[T2, T3], R3 Runnable[T3, T4], R4 Runnable[T4, T5], R5 Runnable[T5, T6], R6 Runnable[T6, T7], R7 Runnable[T7, T8], R8 Runnable[T8, T9], T1, T2, T3, T4, T5, T6, T7, T8, T9 any](r1 R1, r2 R2, r3 R3, r4 R4, r5 R5, r6 R6, r7 R7, r8 R8) Runnable[T1, T9] {
	return chain8[T1, T2, T3, T4, T5, T6, T7, T8, T9]{
		r1: r1,
		r2: r2,
		r3: r3,
		r4: r4,
		r5: r5,
		r6: r6,
		r7: r7,
		r8: r8,
	}
}

type chain9[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any] struct {
	r1 Runnable[T1, T2]
	r2 Runnable[T2, T3]
	r3 Runnable[T3, T4]
	r4 Runnable[T4, T5]
	r5 Runnable[T5, T6]
	r6 Runnable[T6, T7]
	r7 Runnable[T7, T8]
	r8 Runnable[T8, T9]
	r9 Runnable[T9, T10]
}

func (c chain9[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]) Invoke(in T1) (out T10, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	y, err := c.r2.Invoke(x)
	if err != nil {
		return
	}
	z, err := c.r3.Invoke(y)
	if err != nil {
		return
	}
	w, err := c.r4.Invoke(z)
	if err != nil {
		return
	}
	u, err := c.r5.Invoke(w)
	if err != nil {
		return
	}
	v, err := c.r6.Invoke(u)
	if err != nil {
		return
	}
	t, err := c.r7.Invoke(v)
	if err != nil {
		return
	}
	s, err := c.r8.Invoke(t)
	if err != nil {
		return
	}
	return c.r9.Invoke(s)
}

// Chain9 takes 9 Runnable instances and returns a new Runnable instance that chains the nine together.
func Chain9[R1 Runnable[T1, T2], R2 Runnable[T2, T3], R3 Runnable[T3, T4], R4 Runnable[T4, T5], R5 Runnable[T5, T6], R6 Runnable[T6, T7], R7 Runnable[T7, T8], R8 Runnable[T8, T9], R9 Runnable[T9, T10], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](r1 R1, r2 R2, r3 R3, r4 R4, r5 R5, r6 R6, r7 R7, r8 R8, r9 R9) Runnable[T1, T10] {
	return chain9[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]{
		r1: r1,
		r2: r2,
		r3: r3,
		r4: r4,
		r5: r5,
		r6: r6,
		r7: r7,
		r8: r8,
		r9: r9,
	}
}

type chain10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any] struct {
	r1  Runnable[T1, T2]
	r2  Runnable[T2, T3]
	r3  Runnable[T3, T4]
	r4  Runnable[T4, T5]
	r5  Runnable[T5, T6]
	r6  Runnable[T6, T7]
	r7  Runnable[T7, T8]
	r8  Runnable[T8, T9]
	r9  Runnable[T9, T10]
	r10 Runnable[T10, T11]
}

func (c chain10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]) Invoke(in T1) (out T11, err error) {
	x, err := c.r1.Invoke(in)
	if err != nil {
		return
	}
	y, err := c.r2.Invoke(x)
	if err != nil {
		return
	}
	z, err := c.r3.Invoke(y)
	if err != nil {
		return
	}
	w, err := c.r4.Invoke(z)
	if err != nil {
		return
	}
	u, err := c.r5.Invoke(w)
	if err != nil {
		return
	}
	v, err := c.r6.Invoke(u)
	if err != nil {
		return
	}
	t, err := c.r7.Invoke(v)
	if err != nil {
		return
	}
	s, err := c.r8.Invoke(t)
	if err != nil {
		return
	}
	r, err := c.r9.Invoke(s)
	if err != nil {
		return
	}
	return c.r10.Invoke(r)
}

// Chain10 takes 10 Runnable instances and returns a new Runnable instance that chains the ten together.
func Chain10[R1 Runnable[T1, T2], R2 Runnable[T2, T3], R3 Runnable[T3, T4], R4 Runnable[T4, T5], R5 Runnable[T5, T6], R6 Runnable[T6, T7], R7 Runnable[T7, T8], R8 Runnable[T8, T9], R9 Runnable[T9, T10], R10 Runnable[T10, T11], T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11 any](r1 R1, r2 R2, r3 R3, r4 R4, r5 R5, r6 R6, r7 R7, r8 R8, r9 R9, r10 R10) Runnable[T1, T11] {
	return chain10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11]{
		r1:  r1,
		r2:  r2,
		r3:  r3,
		r4:  r4,
		r5:  r5,
		r6:  r6,
		r7:  r7,
		r8:  r8,
		r9:  r9,
		r10: r10,
	}
}
