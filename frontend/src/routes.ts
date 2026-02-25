export const ROUTES = {
    home: "/",
    records: "/records",
    contact: "/contact",
    recordDetail: (recordId: string) => `/records/${recordId}`,
}