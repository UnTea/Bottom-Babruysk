import AlbumInfo from "@/components/AlbumInfo";
import Track from "@/components/Track";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/album/$albumId")({
  component: RouteComponent,
});

class TrackM {
  number: number;
  title: string = "";
  artists: string[] = [];

  constructor(number: number, title: string, artists: string[]) {
    this.number = number;
    this.title = title;
    this.artists = artists;
  }
}

const tracks = [
  new TrackM(1, "On The Rivers Of Mississippi (Intro)", ["Alice Francis"]),
  new TrackM(2, "St. James Ballroom", ["Alice Francis"]),
  new TrackM(3, "Don't Shoot Me (Interlude)", ["Alice Francis"]),
  new TrackM(4, "Shoot Him Down!", ["Alice Francis"]),
  new TrackM(5, "French Affair", ["Alice Francis"]),
  new TrackM(6, "The Funeral (Interlude)", ["Alice Francis"]),
  new TrackM(7, "Gangsterlove", ["Alice Francis"]),
  new TrackM(8, "Sista", ["Alice Francis"]),
  new TrackM(9, "Cakes & Applepies", ["Alice Francis"]),
  new TrackM(10, "Get A Wiggle On", ["Alice Francis"]),
  new TrackM(11, "My Dedication (Interlude)", ["Alice Francis"]),
  new TrackM(12, "I Pimp You", ["Alice Francis"]),
  new TrackM(13, "Chillin' With Niegl (Interlude)", ["Alice Francis"]),
  new TrackM(14, "A Flapper's Diary", ["Alice Francis"]),
  new TrackM(14, "Kiss My Ass", ["Alice Francis"]),
  new TrackM(14, "Please, Love Me Too", ["Alice Francis"]),
  new TrackM(14, "Sandman", ["Alice Francis"]),
];

function RouteComponent() {
  const { albumId } = Route.useParams();
  return (
    <div className="flex w-[1000px] ml-auto mr-auto bg-black">
      <div className="flex flex-col items-center text-center">
        <AlbumInfo
          coverArtLink=""
          title="St. James Ballroom"
          artists={["Alice Francis"]}
        ></AlbumInfo>
      </div>
      <div className="flex-1">
        <p className="text-sm m-2 font-bold text-neutral-200">TRACKS</p>
        {tracks.map((track) => (
          <Track
            artists={track.artists}
            title={track.title}
            number={track.number}
          ></Track>
        ))}
      </div>
    </div>
  );
}
