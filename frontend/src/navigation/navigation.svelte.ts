import type { Route } from "./navigation.types";

export const navigation = $state({
  currentRoute: "servers" as Route,
});

export function navigate(route: Route): void {
  navigation.currentRoute = route;
}
