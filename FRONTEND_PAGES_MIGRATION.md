# Migração de responsabilidades

O frontend segue a separação `Page.svelte` (composição), `store.svelte.ts`
(estado e coordenação), `api.ts` (Wails) e `service.ts` (regras puras).

## Estado atual

- `config`, `images`, `networks` e `stacks`: fronteiras de página migradas.
- `profiles`, `servers`, `containers` e `volumes`: stores e APIs separados.
- `Browse` de stacks está encapsulado na infraestrutura do domínio de stacks.
- Regras puras do backend existem para containers, images, networks, profiles,
  servers e volumes.
- Infraestrutura compartilhada inicial: `internal/shared/events` e
  `internal/shared/filesystem`.

## Validação

O backend é validado com `go test ./...`; o frontend com `npm run check` e
`npm run build` a partir de `frontend/`.

As próximas melhorias estruturais devem preservar os contratos gerados em
`frontend/bindings` e ampliar os testes dos services e repositories.
