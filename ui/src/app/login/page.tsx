import { auth } from '@/auth';
import LoginForm from '@/components/LoginForm';
import { redirect } from 'next/navigation';
import { Suspense } from 'react';

export default async function LoginPage() {
    const session = await auth();

    if (session?.user) {
        redirect('/');
    }

    return (
        <main className="flex items-center justify-center md:h-[calc(100vh-4rem)] bg-zinc-50 font-sans dark:bg-gray-600">
            <Suspense>
                <LoginForm />
            </Suspense>
        </main >
    );
}