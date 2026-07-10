import { mountReactIslands, onDomReady } from "@/lib/reactIslands";

import { ProfilePage } from "..";

onDomReady(() => {
  mountReactIslands("[data-profile-page]", ProfilePage);
});
