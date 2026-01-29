FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

COPY player-service.csproj .
RUN dotnet restore

COPY . .
RUN dotnet publish -c Release -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS runtime
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends curl && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/publish .
COPY --from=build /src/Data ./Data

ENV ASPNETCORE_URLS=http://+:5043
ENV ASPNETCORE_ENVIRONMENT=Production

EXPOSE 5043

ENTRYPOINT ["dotnet", "player-service.dll"]
