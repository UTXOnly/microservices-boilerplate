"""Item service - async CRUD operations."""

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.models.item import Item
from app.schemas.item import ItemCreate, ItemUpdate


class ItemService:
    """Service class for Item operations - reusable across routers and tasks."""

    def __init__(self, db: AsyncSession):
        self.db = db

    async def get_all(self, *, skip: int = 0, limit: int = 100) -> list[Item]:
        """Fetch items with pagination."""
        result = await self.db.execute(
            select(Item).offset(skip).limit(limit).order_by(Item.created_at.desc())
        )
        return list(result.scalars().all())

    async def get_by_id(self, item_id: int) -> Item | None:
        """Fetch a single item by ID."""
        result = await self.db.execute(select(Item).where(Item.id == item_id))
        return result.scalar_one_or_none()

    async def create(self, data: ItemCreate) -> Item:
        """Create a new item."""
        item = Item(**data.model_dump())
        self.db.add(item)
        await self.db.flush()
        await self.db.refresh(item)
        return item

    async def update(self, item: Item, data: ItemUpdate) -> Item:
        """Update an existing item (partial update)."""
        update_data = data.model_dump(exclude_unset=True)
        for key, value in update_data.items():
            setattr(item, key, value)
        await self.db.flush()
        await self.db.refresh(item)
        return item

    async def delete(self, item: Item) -> None:
        """Delete an item."""
        await self.db.delete(item)
        await self.db.flush()
