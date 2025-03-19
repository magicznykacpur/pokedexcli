package pokedex

import (
	"sync"

	"github.com/magicznykacpur/pokedexcli/internal/decoding"
)

type Pokedex struct {
	caughtPokemons map[string]decoding.Pokemon
	mu             sync.Mutex
}

func (p *Pokedex) Catch(pokemon decoding.Pokemon) {
	p.mu.Lock()
	p.caughtPokemons[pokemon.Name] = pokemon
	p.mu.Unlock()
}

func (p *Pokedex) Get(name string) (decoding.Pokemon, bool){
	p.mu.Lock()
	pokemon, ok := p.caughtPokemons[name]
	p.mu.Unlock()
	
	return pokemon, ok
}
func NewPokedex() Pokedex {
	return Pokedex{caughtPokemons: map[string]decoding.Pokemon{}}
}