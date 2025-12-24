import { signOut } from "@/auth";
import { LogIn, LogOut } from "lucide-react";
import { User } from "next-auth";
import Link from "next/link";


function LogoutButton({ user }: { user?: User }) {
    return (
        <form
            action={async () => {
                "use server"
                await signOut({ redirectTo: '/' });
            }}
        >
            <button className="p-2" >
                {user?.email === undefined ?
                    <Link href="/login">
                        <LogIn size={24} />
                    </Link> :
                    <div className="flex items-center gap-2">
                        {user?.email}
                        < LogOut size={24} />
                    </div>
                }
            </button>
        </form>

    )
}


export default LogoutButton;



