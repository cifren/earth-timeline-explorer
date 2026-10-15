import { Timeline } from '@bc/chronology'
import { WorldMap, useGeography } from '@bc/geography'
import { PhotoPanel, useImagery } from '@bc/imagery'
import { useSelection } from '@core/selection/use-selection'

/** Map, timeline and photo panel: the whole v0 experience on one screen. */
export const ExplorerPage = () => {
  const { outlines, isError, zoneName } = useGeography()
  const { zoneId, decade, image, availableZoneIds } = useImagery()
  const { selectZone, selectDecade } = useSelection()

  return (
    <div className="flex h-screen flex-col overflow-hidden">
      <main className="relative min-h-0 flex-1">
        <h1 className="sr-only">Frise Chrono</h1>
        <WorldMap
          outlines={outlines}
          availableZoneIds={availableZoneIds}
          selectedZoneId={zoneId}
          decade={decade}
          marker={image?.position}
          onSelectZone={selectZone}
        />
        {isError && (
          <p role="alert" className="absolute top-20 left-3 rounded bg-card/85 px-2 py-1 text-xs">
            The map is temporarily unavailable.
          </p>
        )}
        <PhotoPanel zoneName={zoneId === null ? undefined : zoneName(zoneId)} />
      </main>
      <Timeline decade={decade} onSelectDecade={selectDecade} />
    </div>
  )
}
