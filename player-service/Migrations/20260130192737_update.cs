using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace player_service.Migrations
{
    /// <inheritdoc />
    public partial class update : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "AquiredAt",
                table: "Inventories",
                newName: "AcquiredAt");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "AcquiredAt",
                table: "Inventories",
                newName: "AquiredAt");
        }
    }
}
