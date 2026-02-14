"""Items CRUD API."""

from fastapi import APIRouter, Depends, HTTPException

from app.database import get_db
from app.schemas.item import ItemCreate, ItemRead, ItemUpdate
from app.services.item_service import ItemService
from sqlalchemy.ext.asyncio import AsyncSession

router = APIRouter(prefix="/items", tags=["Items"])


@router.get("", response_model=list[ItemRead])
async def list_items(
    skip: int = 0,
    limit: int = 100,
    db: AsyncSession = Depends(get_db),
) -> list[ItemRead]:
    """
    List items with pagination.
    - **skip**: Number of records to skip (default: 0)
    - **limit**: Maximum records to return (default: 100)
    """
    service = ItemService(db)
    items = await service.get_all(skip=skip, limit=limit)
    return items


@router.get("/{item_id}", response_model=ItemRead)
async def get_item(
    item_id: int,
    db: AsyncSession = Depends(get_db),
) -> ItemRead:
    """Get a single item by ID."""
    service = ItemService(db)
    item = await service.get_by_id(item_id)
    if not item:
        raise HTTPException(status_code=404, detail="Item not found")
    return item


@router.post("", response_model=ItemRead, status_code=201)
async def create_item(
    data: ItemCreate,
    db: AsyncSession = Depends(get_db),
) -> ItemRead:
    """Create a new item."""
    service = ItemService(db)
    return await service.create(data)


@router.patch("/{item_id}", response_model=ItemRead)
async def update_item(
    item_id: int,
    data: ItemUpdate,
    db: AsyncSession = Depends(get_db),
) -> ItemRead:
    """Partially update an item."""
    service = ItemService(db)
    item = await service.get_by_id(item_id)
    if not item:
        raise HTTPException(status_code=404, detail="Item not found")
    return await service.update(item, data)


@router.delete("/{item_id}", status_code=204)
async def delete_item(
    item_id: int,
    db: AsyncSession = Depends(get_db),
) -> None:
    """Delete an item."""
    service = ItemService(db)
    item = await service.get_by_id(item_id)
    if not item:
        raise HTTPException(status_code=404, detail="Item not found")
    await service.delete(item)
