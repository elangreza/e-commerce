import { GetDetailsProducts, ListProductResponse } from "./product";


export interface DataResponse {
    data: ListProductResponse | GetDetailsProducts;
}