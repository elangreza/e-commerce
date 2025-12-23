import { Product } from "@/types/product";
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { immer } from 'zustand/middleware/immer';

type Cart = {
    UserID: string;
    Items: CartItem[];
    totalPrice: number;
    isLoading: boolean;
    errorMessage: string;
}

export type CartItem = {
    ProductID: string;
    Quantity: number;
}

type Actions = {
    addCartItem: (productID: CartItem['ProductID'], quantity: number) => void
    addQuantityInCart: (productID: CartItem['ProductID'], quantity: number) => void
    removeCartItem: (productID: CartItem['ProductID']) => void
    updateUserID: (s: string) => void
    calculateTotalPrice: (products: Product[]) => void
    setErrorMessage: (msg: string) => void
}

export const useCartStore = create<Cart & Actions>()(
    persist(
        immer((set) => ({
            UserID: "",
            Items: [],
            totalPrice: 0,
            isLoading: false,
            errorMessage: "",
            addCartItem: (productID: CartItem['ProductID'], quantity: number) => {
                if (quantity === 0) {
                    return
                }

                set((state) => {
                    state.isLoading = true
                    state.errorMessage = ""
                    var item = state.Items.find((item) => item.ProductID === productID)
                    if (item) {
                        item.Quantity = quantity;
                        state.isLoading = false
                        return
                    }

                    state.Items.push({
                        ProductID: productID,
                        Quantity: quantity,
                    })
                })

                setTimeout(() => {
                    set((state) => {
                        state.isLoading = false
                    })
                }, 1000)
            },
            addQuantityInCart: (productID: CartItem['ProductID'], quantity: number) => {
                if (quantity === 0) {
                    return
                }

                set((state) => {
                    state.isLoading = true
                    state.errorMessage = ""
                    var item = state.Items.find((item) => item.ProductID === productID)
                    if (item) {
                        item.Quantity += quantity;
                    }
                })

                setTimeout(() => {
                    set((state) => {
                        state.isLoading = false
                    })
                }, 1000)
            },
            removeCartItem: (productID: string) => set((state) => {
                state.isLoading = true
                state.errorMessage = ""
                state.Items = state.Items.filter((item) => item.ProductID !== productID);
                if (state.Items.length === 0) {
                    state.totalPrice = 0
                }
                state.isLoading = false
            }),
            updateUserID: (s: string) => set((state) => {
                state.isLoading = true
                state.errorMessage = ""
                state.UserID = s
                state.isLoading = false
            }),
            calculateTotalPrice: (products: Product[]) => {
                set((state) => {
                    state.isLoading = true
                    state.errorMessage = ""
                    var total = 0
                    for (const item of state.Items) {
                        const product = products.find((p) => p.id === item.ProductID)
                        if (product) {
                            total += product.price.units * item.Quantity
                        }
                    }
                    state.totalPrice = total
                })
                setTimeout(() => {
                    set((state) => {
                        state.isLoading = false
                    })
                }, 1000)
            },
            setErrorMessage: (msg: string) => {
                set((state) => {
                    state.errorMessage = msg
                })

                setTimeout(() => {
                    set((state) => {
                        state.errorMessage = ""
                    })
                }, 1000)
            }
        })), {
        name: "cart",
    })
)