
interface HeaderTextProps {
    name: string;
}

function HeaderText({ name }: HeaderTextProps) {
    return (
        <h1 className="text-2xl font-bold text-white" >{name}</h1>
    )
}


export default HeaderText;