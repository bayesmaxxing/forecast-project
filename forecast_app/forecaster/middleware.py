class CorsMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        response = self.get_response(request)
        response["Access-Control-Allow-Origin"] = "*"  # Allow any domain
        response["Access-Control-Allow-Methods"] = "GET, OPTIONS, PATCH, POST, PUT, DELETE"
        response["Access-Control-Allow-Headers"] = "X-Requested-With, Content-Type"
        return response