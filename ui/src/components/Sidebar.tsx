import { Link } from "@tanstack/react-router";
import SidebarButton from "./SidebarButton";

export default function Sidebar() {
  return (
    <header className="p-2 flex gap-2 bg-neutral-950 justify-between">
      <nav className="flex flex-col items-center">
        <Link to="/">
          <SidebarButton text="Home"></SidebarButton>
        </Link>
        <Link to="/">
          <SidebarButton text="Albums"></SidebarButton>
        </Link>
      </nav>
    </header>
  );
}
