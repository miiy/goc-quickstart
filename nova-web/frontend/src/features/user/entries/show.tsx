import { mountReactIslands, onDomReady } from "@/lib/reactIslands";

import { UserShowPage } from "..";

onDomReady(() => {
  mountReactIslands("[data-user-show]", UserShowPage);
});
