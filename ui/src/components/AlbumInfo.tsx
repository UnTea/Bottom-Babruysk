type AlbumInfoProps = {
  title: string;
  coverArtLink: string;
  artists: string[];
};

export default function AlbumInfo(props: AlbumInfoProps) {
  return (
    <div className="flex flex-col bg-black text-white items-center max-w-[400px]">
      <div className="bg-neutral-500 m-4">
        <img
          className="rounded-lg"
          src={props.coverArtLink}
          width={300}
          height={300}
        />
      </div>

      <h1 className="text-3xl font-bold">{props.title}</h1>
      <p className="text-xl font-semibold">{props.artists} • 2013</p>
    </div>
  );
}
