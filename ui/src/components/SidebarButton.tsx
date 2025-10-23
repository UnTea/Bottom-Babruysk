type SidebarButtonProps = {
  text: string;
};

export default function SidebarButton(props: SidebarButtonProps) {
  return (
    <div className="flex flex-col items-center">
      <p>ICON</p>
      <p>{props.text}</p>
    </div>
  );
}
