'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import Logo from '@/public/Logo.png'

export default function LoginPage() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [showPassword, setShowPassword] = useState(false)

    function handleSubmit(e: React.FormEvent) {
        e.preventDefault()
        console.log({ email, password })
    }

    return (
        <div className="flex min-h-screen items-center justify-center bg-[#FAF7F2] px-4">
            <div className="w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl">
                
                {/* LOGO + TÍTULO */}
                <div className="mb-6 flex flex-col items-center gap-2">
                    <Image
                        src={Logo}
                        alt="Begin Logo"
                        width={100}
                        height={100}
                        priority
                    />

                    <h1 className="text-2xl font-bold text-[#3A2A23]">
                        Entrar
                    </h1>
                    <p className="text-sm text-[#7A5C4D]">
                        Área do Professor
                    </p>
                </div>

                {/* FORM */}
                <form onSubmit={handleSubmit} className="space-y-4">
                    
                    {/* EMAIL */}
                    <div>
                        <label className="mb-1 block text-sm text-[#3A2A23]">
                            Email
                        </label>
                        <input
                            type="email"
                            placeholder="seu@email.com"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            className="w-full rounded-lg border border-[#E5E5E5] px-3 py-2 outline-none focus:border-[#0B3CDE]"
                        />
                    </div>

                    {/* PASSWORD */}
                    <div>
                        <label className="mb-1 block text-sm text-[#3A2A23]">
                            Password
                        </label>

                        <div className="relative">
                            <input
                                type={showPassword ? 'text' : 'password'}
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                className="w-full rounded-lg border border-[#E5E5E5] px-3 py-2 pr-10 outline-none focus:border-[#0B3CDE]"
                            />

                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400"
                            >
                                👁️
                            </button>
                        </div>
                    </div>

                    {/* BOTÃO */}
                    <button
                        type="submit"
                        className="mt-2 w-full rounded-xl bg-[#0B3CDE] py-3 font-semibold text-white transition hover:brightness-110"
                    >
                        Entrar
                    </button>
                </form>

                {/* LINK */}
                <p className="mt-4 text-center text-sm text-[#0B3CDE]">
                    Não tem conta?{' '}
                    <Link href="/cadastro" className="font-semibold underline">
                        Criar conta
                    </Link>
                </p>
            </div>
        </div>
    )
}