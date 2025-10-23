import { Outlet, createRootRoute } from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { TanstackDevtools } from "@tanstack/react-devtools";

import Player from "@/components/Player";
import Sidebar from "@/components/Sidebar";

export const Route = createRootRoute({
  component: () => (
    <>
      <div className="flex flex-row">
        <Sidebar />
        <Outlet />
      </div>
      <Player />
      <TanstackDevtools
        config={{
          position: "bottom-left",
        }}
        plugins={[
          {
            name: "Tanstack Router",
            render: <TanStackRouterDevtoolsPanel />,
          },
        ]}
      />
    </>
  ),
});
