import { Link, Outlet } from 'react-router'

/** Frame shared by every page: the logo stays at the top left, above the page content. */
export const AppLayout = () => (
  <>
    <Outlet />
    <Link
      to="/"
      className="fixed top-3 left-3 z-10 flex items-center gap-2.5 rounded-xl bg-card/95 py-1.5 pr-4 pl-1.5 shadow-lg ring-2 ring-night/80 transition-transform hover:scale-105"
    >
      <img src="/favicon.svg" alt="" className="size-11 rounded-lg" />
      <span className="font-serif text-xl font-semibold tracking-wide">Frise Chrono</span>
    </Link>
  </>
)
