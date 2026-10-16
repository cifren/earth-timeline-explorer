import { LoadingImage } from '@shared/ui/loading-image'

import type { Image, Link } from '../domain/image'

type PhotoInfoProps = {
  image: Image
  zoneName: string | undefined
}

type PhotoDetailsProps = PhotoInfoProps & {
  onEnlarge: () => void
}

export const PhotoDetails = ({ image, zoneName, onEnlarge }: PhotoDetailsProps) => (
  <article className="space-y-3 text-sm">
    <button type="button" className="block w-full cursor-zoom-in" title="Click to enlarge" onClick={onEnlarge}>
      <LoadingImage
        src={image.thumbnailUrl}
        alt={image.title}
        className="aspect-4/3 w-full rounded bg-muted object-cover"
      />
    </button>
    <PhotoInfo image={image} zoneName={zoneName} />
  </article>
)

/** Title, date, credits and links of a photo. */
export const PhotoInfo = ({ image, zoneName }: PhotoInfoProps) => (
  <div className="space-y-3 text-sm">
    <div>
      <h3 className="line-clamp-3 text-base font-semibold" title={image.title}>
        {image.title}
      </h3>
      <p className="text-[15px]">{image.dateLabel}</p>
    </div>
    <dl className="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-[13px]">
      <dt className="text-muted-foreground">Zone</dt>
      <dd>{zoneName}</dd>
      {image.author && (
        <>
          <dt className="text-muted-foreground">Author</dt>
          <dd>{image.author}</dd>
        </>
      )}
      <dt className="text-muted-foreground">License</dt>
      <dd>{image.license}</dd>
    </dl>
    <p className="inline-block rounded-full bg-secondary px-2 py-0.5 text-xs text-secondary-foreground">
      Source: {image.authority}
      {image.originalHolder && ` · original: ${image.originalHolder}`}
    </p>
    <nav aria-label="Links" className="space-y-0.5 text-[13px]">
      <a className="block text-link hover:underline" href={image.sourceUrl} target="_blank" rel="noreferrer">
        View on {image.authority}
      </a>
      <LinkGroup title="On a map" links={image.mapLinks} />
      <LinkGroup title="Learn more" links={image.learnMore} />
    </nav>
  </div>
)

const LinkGroup = ({ title, links }: { title: string; links: Link[] }) =>
  links.length > 0 && (
    <>
      <h4 className="pt-2 font-semibold text-muted-foreground">{title}</h4>
      {links.map((link) => (
        <a key={link.url} className="block text-link hover:underline" href={link.url} target="_blank" rel="noreferrer">
          {link.label}
        </a>
      ))}
    </>
  )
