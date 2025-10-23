type TrackProps = {
  number: number;
  artists: string[];
  title: string;
};

export default function Track(props: TrackProps) {
  return (
    <div className="flex flex-row  p-1  bg-black text-white hover:bg-neutral-800 rounded-lg transition ease-in-out">
      <div className="flex items-center m-2 min-w-6 text-neutral-400">
        <p>{props.number}</p>
      </div>
      <div className="flex flex-col flex-1">
        <p>{props.title}</p>
        <p className="text-sm text-neutral-300">{props.artists}</p>
      </div>
      <div className="flex items-center m-2">
        <p>Lk</p>
      </div>
      <div className="flex items-center m-2">
        <p>4:20</p>
      </div>
    </div>
  );
}
