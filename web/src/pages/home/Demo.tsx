import TextTransition, { presets } from 'react-text-transition';

export type Variables = {
  key: string,
  value: string,
  encoded: string,
}

interface Props {
  content: Variables[];
  encrypted: boolean
}

export function Demo({ content, encrypted }: Props) {
  return (
    <ul className="bg-dark-100 flex h-full w-full flex-col gap-2 rounded-lg p-5 lg:w-96">
      {content.map(({ key, value, encoded }) => (
        <li key={key} className="text-main flex w-96 flex-row font-semibold">
          <p className="font-normal text-white">
            {key}=
          </p>
          <TextTransition springConfig={presets.gentle}>"{encrypted ? encoded : value}"</TextTransition>
        </li>
      ))}
    </ul>
  );
}
