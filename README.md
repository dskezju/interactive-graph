# interactive-graph

A visualization and editing tool for knowledge graph.

## Project setup
```
yarn install
cp /src/config.example.ts /src/config.ts 
cp /backend/config.go.example /backend/config.go 
```

### Compiles and hot-reloads for development
```
yarn serve
cd backend && go run main.go config.go
```

### Compiles and minifies for production
```
yarn build
```

### Run your unit tests
```
yarn test:unit
```

### Run your end-to-end tests
```
yarn test:e2e
```

### Lints and fixes files
```
yarn lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
