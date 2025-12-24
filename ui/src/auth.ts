// import type { User } from '@/app/lib/definitions';
import NextAuth from 'next-auth';
import Credentials from 'next-auth/providers/credentials';
// import postgres from 'postgres';
import { authConfig } from '@/auth.config';
import { z } from 'zod';

// const sql = postgres(process.env.POSTGRES_URL!, { ssl: 'require' });


export type User = {
    // id: string;
    email: string;
    token: string;
}


export const { auth, signIn, signOut } = NextAuth({
    ...authConfig,
    callbacks: {
        async jwt({ token, user }) {
            if (user) {
                token.user = user;
            }
            return token;
        },
        async session({ session, token }) {
            if (token.user) {
                // @ts-ignore
                session.user = token.user as User;
            }
            return session;
        },
    },
    providers: [
        Credentials({
            async authorize(credentials) {
                const parsedCredentials = z
                    .object({ email: z.string().email(), password: z.string().min(4) })
                    .safeParse(credentials);

                if (parsedCredentials.success) {
                    const { email, password } = parsedCredentials.data;
                    const token = await fetch('http://localhost:8080/auth/login', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ email, password }),
                    });
                    if (!token) return null;
                    const data = await token.json();
                    console.log(data);
                    var res: User = {
                        email: email,
                        token: data.token,
                    }

                    return res;
                    // if (data.success) {
                    //     return data.token;
                    // }
                }

                return null;
            },
        }),
    ],
});