package main

import "github.com/pedregon/scaffold"
import "fmt"
import "context"

type (
	Greeter struct {
		Greetings []string
	}
	Fareweller struct {
		Salutations []string
	}
	Spec interface {
		OnRun() *scaffold.Hook[*Greeter]
		OnClose() *scaffold.TaggedHook[*Fareweller]
	}
	App struct {
		greeter    *scaffold.Hook[*Greeter]
		fareweller *scaffold.TaggedHook[*Fareweller]
	}
)

func (greeter *Greeter) Greet() {
	for _, greeting := range greeter.Greetings {
		fmt.Println(greeting)
	}
}

func (fareweller *Fareweller) Farewell() {
	for _, farewell := range fareweller.Salutations {
		fmt.Println(farewell)
	}
}

func (*Fareweller) Tags() []string {
	return []string{"romantic"}
}

func (app *App) OnRun() *scaffold.Hook[*Greeter] {
	return app.greeter
}

func (app *App) OnClose() *scaffold.TaggedHook[*Fareweller] {
	return app.fareweller
}

type (
	EnPlugin[T Spec] struct{}
)

func (plg EnPlugin[T]) String() string {
	return "en_plugin-1.0.0"
}

func (plg EnPlugin[T]) Mount(c scaffold.Context[T]) error {
	c.App.OnRun().Add(func(greeter *Greeter) error {
		greeter.Greetings = append(greeter.Greetings, "hello")
		return nil
	})
	if err := c.Lazy("es_plugin-1.0.0"); err != nil {
		return err
	}
	c.App.OnClose().Add(func(f *Fareweller) error {
		f.Salutations = append(f.Salutations, "goodbye")
		return nil
	})
	return nil
}

type (
	EsPlugin[T Spec] struct{}
)

func (plg EsPlugin[T]) String() string {
	return "es_plugin-1.0.0"
}

func (plg EsPlugin[T]) Mount(c scaffold.Context[T]) error {
	c.App.OnRun().Add(func(greeter *Greeter) error {
		greeter.Greetings = append(greeter.Greetings, "hola")
		return nil
	})
	c.App.OnClose().Add(func(f *Fareweller) error {
		f.Salutations = append(f.Salutations, "adios")
		return nil
	})
	return nil
}

func main() {
	app := &App{
		greeter:    &scaffold.Hook[*Greeter]{},
		fareweller: scaffold.Tag(&scaffold.Hook[*Fareweller]{}),
	}
	scaff := scaffold.New[*App]()
	en := EnPlugin[*App]{}
	es := EsPlugin[*App]{}
	scaff.Register(en)
	scaff.Register(es)
	scaff.Load(context.Background(), app)
	greeter := &Greeter{}
	fareweller := &Fareweller{}
	app.greeter.Trigger(greeter)
	app.fareweller.Trigger(fareweller)
	fmt.Println("[*] Running...")
	greeter.Greet()
	fmt.Println("\n[!] Closing...")
	fareweller.Farewell() // notice the order because of the dependency
	fmt.Println("\n[+] Dependency Graph:")
	for _, dep := range scaffold.Graph(scaff) {
		fmt.Println(dep)
	}
}