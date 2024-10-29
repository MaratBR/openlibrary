import { NavLink } from "react-router-dom";
import { TagDto } from "../api";
import clsx from "clsx";

export type TagProps = {
  tag: TagDto;
};

export default function Tag({ tag }: TagProps) {
  tag = {
    ...tag,
    isAdult: true,
  };

  return (
    <NavLink
      data-tag-type={tag.category}
      data-adult={tag.isAdult}
      data-defined={tag.isDefined}
      className={clsx(
        {
          "!border-red-900/50": tag.isAdult,
        },
        "text-sm bg-muted rounded-sm active:outline outline-3 outline-primary outline-offset-[-2px] inline items-center overflow-hidden border-2 border-muted" +
          ""
      )}
      to={`/tag/${encodeURIComponent(tag.name)}`}
    >
      <span className="mx-0.5 inline whitespace-nowrap">{tag.name}</span>
      {tag.isAdult && <span className="bg-red-900/50 px-1">18+</span>}
    </NavLink>
  );
}
