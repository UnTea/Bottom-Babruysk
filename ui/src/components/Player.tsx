export default function Player() {
  return (
    <div className="sticky bottom-0">
      <div className="absolute top-0 bg-white h-1.5 w-[30%] z-0"></div>
      <div className="bg-black/40 backdrop-blur-[28px] rounded-lg z-10">
        <div className="bg-white w-[30%] h-1"></div>
        <div className="flex items-center">
          <div className="flex items-center gap-4 p-5">
            <img src="/art" width={48} height={48} />
            <div>
              <p>
                <a className="hover:underline" href="#">
                  St. James Ballroom
                </a>
              </p>
              <p>
                <a className="hover:underline" href="#">
                  Alice Francis
                </a>
              </p>
            </div>
          </div>
          <div className="flex-1">ab</div>
          <div className="p-5">Volume</div>
        </div>
      </div>
    </div>
  );
}
